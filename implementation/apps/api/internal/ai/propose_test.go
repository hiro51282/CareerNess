package ai

import (
	"testing"

	"careerness/api/internal/extraction"
	"careerness/api/internal/patch"
)

// TestPropose_ProducesCanonicalFact は mock の出力が生ブロブではなく
// 正規 YAMLFact ベースの upsert_fact パッチであることを検証する（D1）。
func TestPropose_ProducesCanonicalFact(t *testing.T) {
	res := Propose(&ProposeRequest{
		SessionID:   "sess-test",
		WorkspaceID: "my-vault",
		Message:     "2023年にABC社で決済基盤をGoへ移行しCIを短縮した",
	})

	if res.Reply == "" {
		t.Error("reply が空")
	}
	if res.Patch == nil {
		t.Fatal("patch が nil")
	}
	// リクエストの workspace_id が優先されること
	if res.Patch.WorkspaceID != "my-vault" {
		t.Errorf("workspace_id = %q, want my-vault", res.Patch.WorkspaceID)
	}
	if len(res.Patch.Operations) != 1 {
		t.Fatalf("operations = %d, want 1", len(res.Patch.Operations))
	}

	op := res.Patch.Operations[0]
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
	if err := patch.Validate(res.Patch); err != nil {
		t.Fatalf("Propose が不正な patch を生成: %v", err)
	}
}

// TestPropose_SummaryTruncation は長い発言が summary 用に短縮されることを確認する。
func TestPropose_SummaryTruncation(t *testing.T) {
	long := ""
	for range 60 {
		long += "あ"
	}
	res := Propose(&ProposeRequest{SessionID: "s", Message: long})
	fact := res.Patch.Operations[0].Change.After.(*extraction.YAMLFact)
	// 40 文字 + 省略記号
	if []rune(fact.Summary)[40] != '…' {
		t.Errorf("summary が 40 文字で省略されていない: %q", fact.Summary)
	}
	// description は全文を保持
	if fact.Description != long {
		t.Error("description は発言全文を保持すべき")
	}
}
