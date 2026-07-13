package ai

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"careerness/api/internal/extraction"
	"careerness/api/internal/patch"
)

// TestPropose_ProducesCanonicalFact は ExtractionService 経由で得た fact が
// 正規 YAMLFact ベースの upsert_fact パッチとして返ることを検証する（既定 Mock 経路）。
func TestPropose_ProducesCanonicalFact(t *testing.T) {
	res, err := Propose(context.Background(), &ProposeRequest{
		SessionID:   "sess-test",
		WorkspaceID: "my-vault",
		Message:     "2023年にABC社で決済基盤をGoへ移行しCIを短縮した",
	})
	if err != nil {
		t.Fatalf("Propose error: %v", err)
	}

	if res.Reply == "" {
		t.Error("reply が空")
	}
	// Mock は本 PR で単一 fact のまま維持するため patches は 1 件（契約は複数対応）。
	if len(res.Patches) != 1 {
		t.Fatalf("patches = %d, want 1", len(res.Patches))
	}
	p0 := res.Patches[0]
	// リクエストの workspace_id が優先されること
	if p0.WorkspaceID != "my-vault" {
		t.Errorf("workspace_id = %q, want my-vault", p0.WorkspaceID)
	}
	if len(p0.Operations) != 1 {
		t.Fatalf("operations = %d, want 1", len(p0.Operations))
	}

	op := p0.Operations[0]
	if op.Type != patch.OpUpsertFact {
		t.Errorf("op type = %q, want upsert_fact", op.Type)
	}
	if op.Target != "facts/experiences.yaml" {
		t.Errorf("target = %q, want facts/experiences.yaml", op.Target)
	}
	if !op.ReviewRequired {
		t.Error("review_required が false")
	}

	fact, ok := op.Change.After.(*extraction.YAMLFact)
	if !ok {
		t.Fatalf("change.after が *YAMLFact でない: %T", op.Change.After)
	}
	if fact.FactID == "" {
		t.Error("fact_id が空")
	}
	if op.EntityID != fact.FactID {
		t.Errorf("entity_id (%q) と fact_id (%q) が不一致", op.EntityID, fact.FactID)
	}
	if fact.Type != "experience" {
		t.Errorf("type = %q, want experience", fact.Type)
	}
	if fact.Status != "proposed" {
		t.Errorf("status = %q, want proposed", fact.Status)
	}
	if fact.Summary == "" || fact.Description == "" {
		t.Error("summary / description が空")
	}
	// Phase 1: action / decision は空のまま
	if fact.Action != "" || fact.Decision != "" {
		t.Errorf("action/decision は空であるべき: action=%q decision=%q", fact.Action, fact.Decision)
	}

	// docs 準拠の patch であること
	if err := patch.Validate(p0); err != nil {
		t.Fatalf("Propose が不正な patch を生成: %v", err)
	}

	// clarification が一級要素として透過されること（既定 mock は 1 件返す）
	if len(res.Clarifications) != 1 {
		t.Errorf("clarifications = %d, want 1", len(res.Clarifications))
	}
}

// TestPropose_ConversationalTurn は非 fact の発言（疑問形）がエラーにならず、
// patches 空＋会話返信のみで応答することを検証する（B: 自由対話）。
func TestPropose_ConversationalTurn(t *testing.T) {
	res, err := Propose(context.Background(), &ProposeRequest{
		SessionID: "sess-test",
		Message:   "他にどんな情報を書けばいいですか？",
	})
	if err != nil {
		t.Fatalf("疑問形の発言はエラーにならないべき: %v", err)
	}
	if len(res.Patches) != 0 {
		t.Errorf("patches = %d, want 0（非 fact の発言）", len(res.Patches))
	}
	if res.Reply == "" {
		t.Error("会話返信 reply が空")
	}
}

// TestBuildTranscript は会話 transcript 組み立ての純関数を検証する。
func TestBuildTranscript(t *testing.T) {
	// 履歴なし → 最新発言をそのまま（従来挙動）
	if got := buildTranscript(nil, "こんにちは"); got != "こんにちは" {
		t.Errorf("履歴なし = %q, want こんにちは", got)
	}

	// 形式・順序・role マッピング（ai → assistant）・改行の畳み込み
	history := []ChatTurn{
		{Role: "user", Text: "ABC社で\nGo移行をしました"},
		{Role: "ai", Text: "期間を教えてください"},
	}
	got := buildTranscript(history, "2024年4月から12月です")
	want := "user: ABC社で Go移行をしました\nassistant: 期間を教えてください\nuser: 2024年4月から12月です"
	if got != want {
		t.Errorf("transcript = %q, want %q", got, want)
	}

	// 空テキストのターンはスキップ
	got = buildTranscript([]ChatTurn{{Role: "user", Text: "  "}}, "最新")
	if got != "最新" {
		t.Errorf("空履歴のみ = %q, want 最新", got)
	}

	// ターン数上限: 直近 maxHistoryTurns 件のみ
	var many []ChatTurn
	for i := range 20 {
		many = append(many, ChatTurn{Role: "user", Text: fmt.Sprintf("t%02d", i)})
	}
	got = buildTranscript(many, "最新")
	if strings.Contains(got, "t09") || !strings.Contains(got, "t10") {
		t.Errorf("直近 %d 件に制限されるべき: %q", maxHistoryTurns, got)
	}

	// 文字数上限: 古いターンから削られ、直近は残る
	long := strings.Repeat("あ", 3000)
	got = buildTranscript([]ChatTurn{
		{Role: "user", Text: long},
		{Role: "user", Text: long},
		{Role: "user", Text: "直近の発言"},
	}, "最新")
	if !strings.Contains(got, "直近の発言") {
		t.Error("直近ターンは文字数上限後も残るべき")
	}
	if len(got) > maxHistoryChars+len("\nuser: 最新")+16 {
		t.Errorf("文字数上限を超過: %d", len(got))
	}
}

// TestPropose_HistoryKeepsFactClean は履歴付きでも fact が最新発言から
// 組み立てられる（transcript が description を汚さない）ことを検証する（mock 経路）。
func TestPropose_HistoryKeepsFactClean(t *testing.T) {
	latest := "2024年にXYZ社で監視基盤をDatadogへ移行した"
	res, err := Propose(context.Background(), &ProposeRequest{
		SessionID: "sess-test",
		Message:   latest,
		History: []ChatTurn{
			{Role: "user", Text: "こんにちは"},
			{Role: "ai", Text: "キャリアについて教えてください"},
		},
	})
	if err != nil {
		t.Fatalf("Propose error: %v", err)
	}
	if len(res.Patches) != 1 {
		t.Fatalf("patches = %d, want 1", len(res.Patches))
	}
	fact := res.Patches[0].Operations[0].Change.After.(*extraction.YAMLFact)
	if fact.Description != latest {
		t.Errorf("description が transcript で汚染: %q", fact.Description)
	}
}

// TestPropose_SummaryTruncation は長い発言が summary 用に短縮されることを確認する。
func TestPropose_SummaryTruncation(t *testing.T) {
	long := ""
	for range 60 {
		long += "あ"
	}
	res, err := Propose(context.Background(), &ProposeRequest{SessionID: "s", Message: long})
	if err != nil {
		t.Fatalf("Propose error: %v", err)
	}
	fact := res.Patches[0].Operations[0].Change.After.(*extraction.YAMLFact)
	// 40 文字 + 省略記号
	if []rune(fact.Summary)[40] != '…' {
		t.Errorf("summary が 40 文字で省略されていない: %q", fact.Summary)
	}
	// description は全文を保持
	if fact.Description != long {
		t.Error("description は発言全文を保持すべき")
	}
}
