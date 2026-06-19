package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"careerness/api/internal/extraction"
	"gopkg.in/yaml.v3"
)

// newFact は test 用の YAMLFact を生成する。
func newFact(id, summary, company string) *extraction.YAMLFact {
	return &extraction.YAMLFact{
		FactID:       id,
		Type:         "experience",
		Status:       "proposed",
		Summary:      summary,
		Description:  "Test description",
		Confidence:   "high",
		Source:       "conversation",
		CreatedAt:    "2026-05-24T10:00:00Z",
		Tags:         []string{"backend"},
		Company:      company,
		Period:       "2022-01 to 2023-01",
		SourceDetail: "Test fact",
	}
}

// upsertPatch は Model A（patch.Patch）形式で upsert_fact パッチを組み立てる。
// change.after に fact 本体を入れる（docs/implementation/ai/ai-patch-format.md）。
func upsertPatch(patchID, target string, fact *extraction.YAMLFact) *Patch {
	return &Patch{
		PatchID:     patchID,
		WorkspaceID: "test-workspace",
		SessionID:   "sess_test",
		CreatedAt:   "2026-05-24T10:00:00Z",
		CreatedBy:   "ai",
		Kind:        "workspace_patch",
		Summary:     "Test patch",
		Status:      StatusProposed,
		Operations: []Operation{
			{
				OpID:            "op-001",
				Type:            OpUpsertFact,
				Target:          target,
				EntityID:        fact.FactID,
				Change:          ChangeRecord{Before: nil, After: fact},
				Rationale:       "test rationale",
				Confidence:      ConfidenceHigh,
				FactStatusAfter: fact.Status,
				ReviewRequired:  true,
			},
		},
	}
}

func TestApplyPatch_UpsertFact(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "careervault-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)
	p := upsertPatch("patch-test-001", "facts/experiences.yaml", newFact("fact-proj-test", "Test project", "Test Corp"))

	result, err := applier.ApplyPatch(p)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("Patch apply failed: %v", result.FailedOps)
	}
	if result.AppliedCount != 1 {
		t.Errorf("AppliedCount = %d, want 1", result.AppliedCount)
	}

	targetPath := filepath.Join(tempDir, "facts", "experiences.yaml")
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("File not created at %s: %v", targetPath, err)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	var facts []*extraction.YAMLFact
	if err := yaml.Unmarshal(content, &facts); err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}
	if len(facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(facts))
	}
	if facts[0].FactID != "fact-proj-test" {
		t.Errorf("FactID = %q, want fact-proj-test", facts[0].FactID)
	}
}

func TestApplyPatch_UpdateExisting(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "careervault-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)

	// 既存ファイルに 1 fact を置く
	existingFact := newFact("fact-proj-existing", "Existing project", "Old Corp")
	factsDir := filepath.Join(tempDir, "facts")
	os.MkdirAll(factsDir, 0755)
	factsFile := filepath.Join(factsDir, "experiences.yaml")
	initialContent, _ := yaml.Marshal([]*extraction.YAMLFact{existingFact})
	os.WriteFile(factsFile, initialContent, 0644)

	// 新 fact を upsert
	p := upsertPatch("patch-test-002", "facts/experiences.yaml", newFact("fact-proj-new", "New project", "New Corp"))
	result, err := applier.ApplyPatch(p)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("Patch apply failed: %v", result.FailedOps)
	}

	content, _ := os.ReadFile(factsFile)
	var facts []*extraction.YAMLFact
	yaml.Unmarshal(content, &facts)

	if len(facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(facts))
	}
	factIDs := make(map[string]bool)
	for _, f := range facts {
		factIDs[f.FactID] = true
	}
	if !factIDs["fact-proj-existing"] {
		t.Error("Existing fact should still be present")
	}
	if !factIDs["fact-proj-new"] {
		t.Error("New fact should be added")
	}
}

func TestApplyPatch_InvalidStatus(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)
	// 構造は valid だが status が不正なケースを検証する。
	p := upsertPatch("patch-bad-status", "facts/experiences.yaml", newFact("fact-proj-bad", "Bad", "Bad Corp"))
	p.Status = "invalid_status"

	if _, err := applier.ApplyPatch(p); err == nil {
		t.Error("ApplyPatch should fail for invalid status")
	}
}

func TestApplyPatch_NoOperations(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)
	p := &Patch{
		PatchID:     "patch-empty",
		WorkspaceID: "test-workspace",
		SessionID:   "sess_test",
		Status:      StatusProposed,
		Operations:  []Operation{},
	}

	if _, err := applier.ApplyPatch(p); err == nil {
		t.Error("ApplyPatch should fail for patch with no operations")
	}
}

// TestApplyPatch_RejectsTraversalTarget は op.Target にトラバーサルを含む patch が
// 入口の Validate で弾かれ、FS へ到達しないことを検証する（ADR-006）。
func TestApplyPatch_RejectsTraversalTarget(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)
	p := upsertPatch("patch-evil", "../../etc/evil.yaml", newFact("fact-proj-evil", "Evil", "Evil Corp"))

	if _, err := applier.ApplyPatch(p); err == nil {
		t.Fatal("traversal target を含む patch は拒否されるべき")
	}
}

// TestApplyPatch_RejectsSymlinkEscape は ".." を含まないため Validate は通るが、
// root 内 symlink が外部を指す target を ResolveWithin が apply 時に遮断することを検証する。
func TestApplyPatch_RejectsSymlinkEscape(t *testing.T) {
	root, _ := os.MkdirTemp("", "careervault-root-")
	defer os.RemoveAll(root)
	outside, _ := os.MkdirTemp("", "careervault-out-")
	defer os.RemoveAll(outside)

	if err := os.Symlink(outside, filepath.Join(root, "out")); err != nil {
		t.Skipf("symlink を作成できない環境: %v", err)
	}

	applier := NewApplier(root)
	p := upsertPatch("patch-symlink", "out/evil.yaml", newFact("fact-proj-evil", "Evil", "Evil Corp"))

	if _, err := applier.ApplyPatch(p); err == nil {
		t.Fatal("symlink 経由の root 外書き込みは遮断されるべき")
	}
	if _, err := os.Stat(filepath.Join(outside, "evil.yaml")); err == nil {
		t.Fatal("root 外（symlink 先）にファイルが作成された")
	}
}

// TestApplyPatch_JSONRoundTrip は最大リスク箇所を検証する：
// patch を JSON に marshal → unmarshal すると change.after は
// map[string]interface{} になる。applier がそこから YAMLFact を正しく
// 復元し、配列(tags)・*int(team_size)・map(impact) を保持できることを確認する。
func TestApplyPatch_JSONRoundTrip(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	teamSize := 5
	fact := newFact("fact-roundtrip", "Round trip project", "RT Corp")
	fact.TechStack = []string{"go", "react"}
	fact.TeamSize = &teamSize
	fact.Impact = map[string]interface{}{"summary": "CI 時間を 40% 短縮"}

	p := upsertPatch("patch-rt-001", "facts/experiences.yaml", fact)

	// JSON 往復（クライアント → サーバの実経路を再現）
	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal patch: %v", err)
	}
	var decoded Patch
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal patch: %v", err)
	}

	applier := NewApplier(tempDir)
	result, err := applier.ApplyPatch(&decoded)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("apply failed: %v", result.FailedOps)
	}

	content, _ := os.ReadFile(filepath.Join(tempDir, "facts", "experiences.yaml"))
	var facts []*extraction.YAMLFact
	if err := yaml.Unmarshal(content, &facts); err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	if len(facts) != 1 {
		t.Fatalf("expected 1 fact, got %d", len(facts))
	}
	got := facts[0]
	if got.FactID != "fact-roundtrip" {
		t.Errorf("FactID = %q, want fact-roundtrip", got.FactID)
	}
	if len(got.TechStack) != 2 || got.TechStack[0] != "go" {
		t.Errorf("TechStack not preserved: %v", got.TechStack)
	}
	if got.TeamSize == nil || *got.TeamSize != 5 {
		t.Errorf("TeamSize not preserved: %v", got.TeamSize)
	}
	if got.Impact == nil || got.Impact["summary"] != "CI 時間を 40% 短縮" {
		t.Errorf("Impact not preserved: %v", got.Impact)
	}
}

// TestBuildFactUpsert は builder が docs 準拠の有効なパッチを生成することを検証する。
func TestBuildFactUpsert(t *testing.T) {
	fact := newFact("fact-build-001", "Built project", "Build Corp")
	p := BuildFactUpsert(fact, "sess_build", 0)

	if err := Validate(p); err != nil {
		t.Fatalf("BuildFactUpsert produced invalid patch: %v", err)
	}
	if len(p.Operations) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(p.Operations))
	}
	op := p.Operations[0]
	if op.Type != OpUpsertFact {
		t.Errorf("op type = %q, want upsert_fact", op.Type)
	}
	if op.EntityID != "fact-build-001" {
		t.Errorf("entity_id = %q, want fact-build-001", op.EntityID)
	}
	if op.FactStatusAfter != "proposed" {
		t.Errorf("fact_status_after = %q, want proposed", op.FactStatusAfter)
	}
	if !op.ReviewRequired {
		t.Error("review_required should be true for fact upsert")
	}
	// change.after が fact であること
	gotFact, ok := op.Change.After.(*extraction.YAMLFact)
	if !ok {
		t.Fatalf("change.after is not *YAMLFact: %T", op.Change.After)
	}
	if gotFact.FactID != "fact-build-001" {
		t.Errorf("change.after.fact_id = %q, want fact-build-001", gotFact.FactID)
	}
}

// TestBuildAndApply は builder → applier の結合経路を検証する。
func TestBuildAndApply(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	fact := newFact("fact-e2e-001", "E2E project", "E2E Corp")
	p := BuildFactUpsert(fact, "sess_e2e", 0)

	applier := NewApplier(tempDir)
	result, err := applier.ApplyPatch(p)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("apply failed: %v", result.FailedOps)
	}

	// builder の target は facts/experiences.yaml（type=experience）
	targetPath := filepath.Join(tempDir, "facts", "experiences.yaml")
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

// TestApplyMarkFactStatus は mark_fact_status が fact_status_after を反映することを検証する。
func TestApplyMarkFactStatus(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	// 既存 fact を proposed で配置
	fact := newFact("fact-status-001", "Status project", "Status Corp")
	factsDir := filepath.Join(tempDir, "facts")
	os.MkdirAll(factsDir, 0755)
	factsFile := filepath.Join(factsDir, "experiences.yaml")
	content, _ := yaml.Marshal([]*extraction.YAMLFact{fact})
	os.WriteFile(factsFile, content, 0644)

	p := &Patch{
		PatchID:     "patch-mark-001",
		WorkspaceID: "test-workspace",
		SessionID:   "sess_test",
		Status:      StatusApproved,
		Operations: []Operation{
			{
				OpID:            "op-001",
				Type:            OpMarkFactStatus,
				Target:          "facts/experiences.yaml",
				EntityID:        "fact-status-001",
				Change:          ChangeRecord{Before: "proposed", After: "confirmed"},
				Rationale:       "user approved",
				Confidence:      ConfidenceHigh,
				FactStatusAfter: "confirmed",
			},
		},
	}

	applier := NewApplier(tempDir)
	result, err := applier.ApplyPatch(p)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("apply failed: %v", result.FailedOps)
	}

	updated, _ := os.ReadFile(factsFile)
	var facts []*extraction.YAMLFact
	yaml.Unmarshal(updated, &facts)
	if facts[0].Status != "confirmed" {
		t.Errorf("status = %q, want confirmed", facts[0].Status)
	}
}
