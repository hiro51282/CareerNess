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

	"careerness/api/internal/extraction"
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
// 実 AI が未統合の間は mock として動作するが、出力は正規の YAMLFact 形に揃える
// （docs/implementation/workspace/fact-schema.md, ai-patch-format.md）。
// patch 組み立ては extract 経路と同じ patch.BuildFactUpsert を再利用する。
func Propose(req *ProposeRequest) *ProposeResult {
	now := time.Now().UTC().Format(time.RFC3339)

	// 発言から正規 YAMLFact を組み立てる。type の自動分類や action/decision の
	// 中身埋めは実 AI 統合（Codex）の責務なので、ここでは既定 experience・空のままにする。
	fact := &extraction.YAMLFact{
		FactID:      fmt.Sprintf("fact-proj-%s-%s", slugify(req.Message), shortID()),
		Type:        "experience",
		Status:      "proposed",
		Summary:     summarize(req.Message),
		Description: req.Message,
		Action:      "", // Phase 1: 空（表示側で description にフォールバック）
		Decision:    "",
		Confidence:  "medium",
		Source:      "conversation",
		CreatedAt:   now,
		Company:     "未確認",
		Tags:        []string{},
	}

	p := patch.BuildFactUpsert(fact, req.SessionID, 0)
	// BuildFactUpsert は workspace_id を既定値にするため、リクエスト指定があれば優先する。
	if req.WorkspaceID != "" {
		p.WorkspaceID = req.WorkspaceID
	}

	reply := buildReply(req.Message, len(req.WorkspaceFiles))

	return &ProposeResult{Reply: reply, Patch: p}
}

// summarize は発言を fact の summary 用に短く整える。
func summarize(message string) string {
	return truncateRunes(strings.TrimSpace(message), 40)
}

// truncateRunes は文字列を rune 単位で n 文字に切り詰め、超過時は省略記号を付ける。
func truncateRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}

func buildReply(message string, fileCount int) string {
	preview := truncateRunes(message, 40)

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
