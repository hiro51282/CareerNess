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

// ChatTurn は会話履歴の 1 ターン。
type ChatTurn struct {
	Role string `json:"role"` // "user" | それ以外は assistant として扱う
	Text string `json:"text"`
}

// ProposeRequest はパッチ提案の生成に必要な入力。
type ProposeRequest struct {
	SessionID      string            `json:"session_id"`
	WorkspaceID    string            `json:"workspace_id"`
	Message        string            `json:"message"`
	History        []ChatTurn        `json:"history,omitempty"` // 直近の会話履歴（任意）
	WorkspaceFiles map[string]string `json:"workspace_files"`
}

// ProposeResult は AI が返す提案結果。
// 1 会話ターンから複数 fact が抽出され得るため、patch は複数返す（1 patch = 1 fact）。
type ProposeResult struct {
	Reply string `json:"reply"`
	// Clarifications は AI の聞き返し（全 fact の clarification_questions の集約）。
	// UI はチップ等の一級要素として表示し、回答は履歴付きの次ターンで fact を enrich する。
	Clarifications []string       `json:"clarifications,omitempty"`
	Patches        []*patch.Patch `json:"patches"`
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

	// 会話履歴があれば transcript（末尾 user: が最新発言）として渡す。
	// これにより clarification への回答が既出 fact の更新（同一 fact_id_hint 再利用）に繋がる。
	conversation := buildTranscript(req.History, req.Message)

	result, err := extraction.NewExtractionService(provider).
		ExtractFromConversation(ctx, conversation, req.SessionID)
	if err != nil {
		return nil, err
	}
	// 抽出 0 件は正常（非 fact の発言）。patches 空＋会話返信のみで応答する。
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

	// AI の会話返信を採用する。空の場合のフォールバック:
	// facts あり → 従来の定型文（提案の確認を促す）、facts なし → キャリアの話への誘導。
	// （0 件時の reply 必須は ValidateAPIResult が保証するため、後者は防御的措置）
	reply := strings.TrimSpace(result.Reply)
	if reply == "" {
		if len(patches) > 0 {
			reply = buildReply(req.Message, len(req.WorkspaceFiles))
		} else {
			reply = "承知しました。担当した仕事・成果・スキルについて教えていただくと、fact 候補を提案できます。"
		}
	}

	return &ProposeResult{Reply: reply, Clarifications: result.Clarifications, Patches: patches}, nil
}

// truncateRunes は文字列を rune 単位で n 文字に切り詰め、超過時は省略記号を付ける。
func truncateRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}

// 履歴の上限。prompt の肥大（レイテンシ/コスト）を抑えるため、
// 直近ターン数と文字数の両方で制限する。
const (
	maxHistoryTurns = 10
	maxHistoryChars = 4000
)

// buildTranscript は会話履歴＋最新発言から transcript を組み立てる純関数。
// 形式は extraction-specification.md の会話拡張に従う:
//
//	user: ...
//	assistant: ...
//	user: <最新発言>
//
// 履歴が無い場合は最新発言をそのまま返す（従来挙動・単一発言）。
// 履歴の各ターンは改行を潰して 1 行に収め、古いターンから上限で切り捨てる。
func buildTranscript(history []ChatTurn, latest string) string {
	if len(history) == 0 {
		return latest
	}
	turns := history
	if len(turns) > maxHistoryTurns {
		turns = turns[len(turns)-maxHistoryTurns:]
	}

	lines := make([]string, 0, len(turns))
	total := 0
	for _, t := range turns {
		text := strings.TrimSpace(strings.ReplaceAll(t.Text, "\n", " "))
		if text == "" {
			continue
		}
		role := "assistant"
		if t.Role == "user" {
			role = "user"
		}
		line := role + ": " + text
		lines = append(lines, line)
		total += len(line) + 1
	}
	// 文字数上限は古いターンから削って満たす（直近の文脈を優先）。
	for len(lines) > 0 && total > maxHistoryChars {
		total -= len(lines[0]) + 1
		lines = lines[1:]
	}
	if len(lines) == 0 {
		return latest
	}
	return strings.Join(lines, "\n") + "\nuser: " + latest
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
