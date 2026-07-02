package extraction

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestMock_ReflectsConversation は mock が固定値でなく会話を反映することを検証する。
func TestMock_ReflectsConversation(t *testing.T) {
	m := NewMockExtractionProvider()
	conv := "2023年にABC社で決済基盤をGoへ移行した"

	res, err := m.ExtractFacts(context.Background(), conv)
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}
	if len(res.ExtractedFacts) != 1 {
		t.Fatalf("facts = %d, want 1", len(res.ExtractedFacts))
	}
	f := res.ExtractedFacts[0]
	if f.Type != "experience" {
		t.Errorf("type = %q, want experience", f.Type)
	}
	if f.Confidence != "low" {
		t.Errorf("confidence = %q, want low（mock 由来）", f.Confidence)
	}
	if f.Description != conv {
		t.Errorf("description は発言全文を保持すべき: %q", f.Description)
	}
	if !strings.HasPrefix(conv, string([]rune(f.Summary)[:1])) {
		t.Errorf("summary が会話を反映していない: %q", f.Summary)
	}
}

// TestMock_PipelinePasses は mock 出力が抽出パイプライン全体を通過することを検証する。
func TestMock_PipelinePasses(t *testing.T) {
	svc := NewExtractionService(NewMockExtractionProvider())
	res, err := svc.ExtractFromConversation(context.Background(), "ABCでGo移行を主導した", "sess-1")
	if err != nil {
		t.Fatalf("ExtractFromConversation error: %v", err)
	}
	if len(res.YAMLFacts) < 1 {
		t.Fatalf("YAMLFacts = %d, want >=1", len(res.YAMLFacts))
	}
	f := res.YAMLFacts[0]

	// experience の必須補完
	if f.Company != "未確認" {
		t.Errorf("company = %q, want 未確認（normalizer 補完）", f.Company)
	}
	if f.Period == "" {
		t.Error("period が空")
	}
	if len(f.Tags) == 0 {
		t.Error("tags が空")
	}
	// fact_id 形式と timestamp 妥当性
	if !strings.HasPrefix(f.FactID, "fact-proj-") {
		t.Errorf("fact_id = %q, want fact-proj- 始まり", f.FactID)
	}
	if _, err := time.Parse(time.RFC3339, f.CreatedAt); err != nil {
		t.Errorf("created_at が RFC3339 でない: %q (%v)", f.CreatedAt, err)
	}
}

// TestMock_EmptyConversation は空会話で抽出 0 件となり、パイプラインがエラーを返すことを検証する。
func TestMock_EmptyConversation(t *testing.T) {
	m := NewMockExtractionProvider()
	res, err := m.ExtractFacts(context.Background(), "   ")
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}
	if len(res.ExtractedFacts) != 0 {
		t.Errorf("空会話の facts = %d, want 0", len(res.ExtractedFacts))
	}

	svc := NewExtractionService(NewMockExtractionProvider())
	if _, err := svc.ExtractFromConversation(context.Background(), "", "sess-1"); err == nil {
		t.Error("空会話はパイプラインでエラーになるべき")
	}
}

// TestMock_Override は SetResponse による上書きが効くことを検証する（テスト互換）。
func TestMock_Override(t *testing.T) {
	m := NewMockExtractionProvider()
	override := &ExtractedFactResult{
		ExtractedFacts: []ExtractedFact{{Type: "skill", FactIDHint: "go", Summary: "Go", Description: "d", Confidence: "high"}},
		ExtractionQuality: ExtractionQuality{OverallConfidence: "high"},
	}
	m.SetResponse(override)

	res, err := m.ExtractFacts(context.Background(), "無視される会話")
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}
	if len(res.ExtractedFacts) != 1 || res.ExtractedFacts[0].Type != "skill" {
		t.Errorf("override が反映されていない: %+v", res.ExtractedFacts)
	}
}

// TestMock_SummaryTruncated は長い発言で summary が省略されることを検証する。
func TestMock_SummaryTruncated(t *testing.T) {
	m := NewMockExtractionProvider()
	long := strings.Repeat("あ", 60)

	res, err := m.ExtractFacts(context.Background(), long)
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}
	f := res.ExtractedFacts[0]
	if []rune(f.Summary)[40] != '…' {
		t.Errorf("summary が 40 文字で省略されていない: %q", f.Summary)
	}
	if f.Description != long {
		t.Error("description は全文を保持すべき")
	}
}
