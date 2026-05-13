// Package ai はワークスペース変更の patch proposal 生成を担う。
// MVP では mock 実装。後から OpenAI / Claude 等の実 AI 呼び出しに差し替える。
package ai

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"careerness/api/internal/patch"
)

// ProposeRequest はパッチ提案の生成に必要な入力。
type ProposeRequest struct {
	SessionID      string            `json:"session_id"`
	WorkspaceID    string            `json:"workspace_id"`
	Message        string            `json:"message"`
	WorkspaceFiles map[string]string `json:"workspace_files"`
}

// ProposeResult は AI が返す提案結果。
type ProposeResult struct {
	Reply string       `json:"reply"`
	Patch *patch.Patch `json:"patch,omitempty"`
}

// Propose はユーザーの発言を受けて patch proposal を生成する。
// 実 AI が未統合の間は mock として動作する。
func Propose(req *ProposeRequest) *ProposeResult {
	now := time.Now()
	patchID := fmt.Sprintf("patch-%s-%s", now.Format("2006-01-02"), shortID())
	entityID := fmt.Sprintf("fact-%s-%s", slugify(req.Message), shortID())

	afterContent := map[string]interface{}{
		"source":       "conversation",
		"raw_text":     req.Message,
		"extracted_at": now.Format(time.RFC3339),
		"status":       "proposed",
	}

	p := &patch.Patch{
		PatchID:     patchID,
		WorkspaceID: req.WorkspaceID,
		SessionID:   req.SessionID,
		CreatedAt:   now.Format(time.RFC3339),
		CreatedBy:   "ai",
		Kind:        "workspace_patch",
		Summary:     "会話からキャリア fact 候補を抽出しました",
		Status:      patch.StatusProposed,
		Operations: []patch.Operation{
			{
				OpID:     "op-001",
				Type:     patch.OpUpsertFact,
				Target:   "facts/inbox.yaml",
				EntityID: entityID,
				Change: patch.ChangeRecord{
					Before: nil,
					After:  afterContent,
				},
				Rationale:       "ユーザーの発言から抽出したキャリア情報です。内容を確認してください。",
				Confidence:      patch.ConfidenceMedium,
				FactStatusAfter: "proposed",
				ReviewRequired:  true,
			},
		},
	}

	reply := buildReply(req.Message, len(req.WorkspaceFiles))

	return &ProposeResult{Reply: reply, Patch: p}
}

func buildReply(message string, fileCount int) string {
	preview := message
	if utf8.RuneCountInString(message) > 40 {
		runes := []rune(message)
		preview = string(runes[:40]) + "…"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("「%s」という発言を受け取りました。\n\n", preview))
	if fileCount > 0 {
		sb.WriteString(fmt.Sprintf("ワークスペースの %d ファイルを参照しました。\n", fileCount))
	}
	sb.WriteString("以下の fact 候補を提案します。内容を確認して承認または却下してください。")
	return sb.String()
}

func slugify(s string) string {
	runes := []rune(s)
	if len(runes) > 12 {
		runes = runes[:12]
	}
	var b strings.Builder
	for _, r := range runes {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			b.WriteRune(r + 32)
		} else {
			b.WriteRune('-')
		}
	}
	result := strings.Trim(b.String(), "-")
	if result == "" {
		return "fact"
	}
	return result
}

func shortID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
