package extraction

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestYAMLFact_ActionDecisionRoundTrip は Phase 1 Minimal Fix で追加した
// action/decision が YAML marshal→unmarshal を通して保持されることを検証する。
func TestYAMLFact_ActionDecisionRoundTrip(t *testing.T) {
	fact := &YAMLFact{
		FactID:      "fact-proj-roundtrip",
		Type:        "achievement",
		Status:      "proposed",
		Summary:     "CI 改善",
		Description: "CI workflow を整理した",
		Action:      "CI workflow を整理した",
		Decision:    "並列化より job 分割とキャッシュ整理を優先した",
		Impact:      map[string]interface{}{"summary": "CI 実行時間を短縮した"},
		Confidence:  "medium",
	}

	raw, err := yaml.Marshal(fact)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got YAMLFact
	if err := yaml.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Action != fact.Action {
		t.Errorf("Action = %q, want %q", got.Action, fact.Action)
	}
	if got.Decision != fact.Decision {
		t.Errorf("Decision = %q, want %q", got.Decision, fact.Decision)
	}
	// impact.summary は既存の Impact map で表現できること
	if got.Impact == nil || got.Impact["summary"] != "CI 実行時間を短縮した" {
		t.Errorf("impact.summary not preserved: %v", got.Impact)
	}
}

// TestYAMLFact_EmptyActionDecisionOmitted は Minimal Fix の核を担保する：
// action/decision が空のとき omitempty により YAML へ出力されないこと。
// これにより既存ファイルが汚れず、Phase 2 拡張時のマイグレーションが不要になる。
func TestYAMLFact_EmptyActionDecisionOmitted(t *testing.T) {
	fact := &YAMLFact{
		FactID:      "fact-proj-empty",
		Type:        "experience",
		Status:      "proposed",
		Summary:     "バックエンド開発",
		Description: "API サーバーを実装した",
		Confidence:  "high",
		// Action / Decision は意図的に空
	}

	raw, err := yaml.Marshal(fact)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(raw)

	if strings.Contains(out, "action:") {
		t.Errorf("empty action should be omitted, got:\n%s", out)
	}
	if strings.Contains(out, "decision:") {
		t.Errorf("empty decision should be omitted, got:\n%s", out)
	}
}
