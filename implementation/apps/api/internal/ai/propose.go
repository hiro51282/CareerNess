// Package ai はワークスペース変更の patch proposal 生成を担う。
// チャット経路の fact 抽出は extraction.ExtractionService に一本化しており、
// ai はその結果を patch proposal（＋チャット返信）へ組み立てる薄い層に徹する。
package ai

import (
	"context"
	"fmt"
	"strings"
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
// 1 会話ターンから複数 fact が抽出され得るため、patch は複数返す（1 patch = 1 fact）。
type ProposeResult struct {
	Reply   string         `json:"reply"`
	Patches []*patch.Patch `json:"patches"`
}

// Propose はユーザーの発言を受けて patch proposal を生成する。
// fact 抽出は extraction.ExtractionService に委譲し（provider は env で選択、既定 Mock）、
// ここでは得られた fact を patch proposal（extract 経路と同じ patch.BuildFactUpsert）と
// チャット返信へ組み立てる。実 AI（Codex CLI）の接続は本 PR では行わない。
//
// PR-A では応答契約（単一 patch）を維持する。複数 fact（patches[]）対応は後続 PR-B。
func Propose(ctx context.Context, req *ProposeRequest) (*ProposeResult, error) {
	provider, err := extraction.NewProviderFromEnv()
	if err != nil {
		return nil, fmt.Errorf("extraction provider の初期化に失敗: %w", err)
	}

	result, err := extraction.NewExtractionService(provider).
		ExtractFromConversation(ctx, req.Message, req.SessionID)
	if err != nil {
		return nil, err
	}
	if len(result.YAMLFacts) == 0 {
		// ExtractFromConversation は 0 件で error を返す契約だが、防御的に確認する。
		return nil, fmt.Errorf("会話から fact を抽出できませんでした")
	}

	// 抽出された全 fact を patch 提案化する（1 patch = 1 fact、extract 経路と同じ組み立て）。
	patches := make([]*patch.Patch, 0, len(result.YAMLFacts))
	for i, fact := range result.YAMLFacts {
		p := patch.BuildFactUpsert(fact, req.SessionID, i)
		// BuildFactUpsert は workspace_id を既定値にするため、リクエスト指定があれば優先する。
		if req.WorkspaceID != "" {
			p.WorkspaceID = req.WorkspaceID
		}
		patches = append(patches, p)
	}

	reply := buildReply(req.Message, len(req.WorkspaceFiles))

	return &ProposeResult{Reply: reply, Patches: patches}, nil
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
