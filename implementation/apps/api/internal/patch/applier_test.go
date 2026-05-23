package patch

import (
	"os"
	"path/filepath"
	"testing"

	"careerness/api/internal/extraction"
	"gopkg.in/yaml.v3"
)

func TestApplyPatch_UpsertFact(t *testing.T) {
	// Create temporary workspace
	tempDir, err := os.MkdirTemp("", "careervault-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)

	// Create a patch with upsert_fact operation
	fact := &extraction.YAMLFact{
		FactID:      "fact-proj-test",
		Type:        "experience",
		Status:      "proposed",
		Summary:     "Test project",
		Description: "Test description",
		Confidence:  "high",
		Source:      "conversation",
		CreatedAt:   "2026-05-24T10:00:00Z",
		Tags:        []string{"backend"},
		Company:     "Test Corp",
		Period:      "2022-01 to 2023-01",
		SourceDetail: "Test fact",
	}

	patch := &extraction.Patch{
		PatchID:     "patch-test-001",
		Kind:        "fact_upsert",
		Status:      "proposed",
		CreatedAt:   "2026-05-24T10:00:00Z",
		CreatedBy:   "ai",
		WorkspaceID: "test-workspace",
		SessionID:   "sess_test",
		Summary:     "Test patch",
		Operations: []extraction.Operation{
			{
				OpID:    "op-001",
				Type:    "upsert_fact",
				Target:  "facts/experiences.yaml",
				FactID:  fact.FactID,
				NewFact: fact,
			},
		},
	}

	// Apply patch
	result, err := applier.ApplyPatch(patch)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}

	if !result.Success() {
		t.Fatalf("Patch apply failed: %v", result.FailedOps)
	}

	if result.AppliedCount != 1 {
		t.Errorf("AppliedCount = %d, want 1", result.AppliedCount)
	}

	// Verify file was created and contains fact
	targetPath := filepath.Join(tempDir, "facts", "experiences.yaml")
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("File not created at %s: %v", targetPath, err)
	}

	// Read and verify content
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

	// Create initial file with one fact
	existingFact := &extraction.YAMLFact{
		FactID:      "fact-proj-existing",
		Type:        "experience",
		Status:      "proposed",
		Summary:     "Existing project",
		Description: "Existing description",
		Confidence:  "high",
		Source:      "manual",
		CreatedAt:   "2026-05-20T00:00:00Z",
		Tags:        []string{"backend"},
		Company:     "Old Corp",
		Period:      "2020-01 to 2021-01",
		SourceDetail: "Manual entry",
	}

	factsDir := filepath.Join(tempDir, "facts")
	os.MkdirAll(factsDir, 0755)
	factsFile := filepath.Join(factsDir, "experiences.yaml")

	initialContent, _ := yaml.Marshal([]*extraction.YAMLFact{existingFact})
	os.WriteFile(factsFile, initialContent, 0644)

	// Now apply a patch with a new fact
	newFact := &extraction.YAMLFact{
		FactID:      "fact-proj-new",
		Type:        "experience",
		Status:      "proposed",
		Summary:     "New project",
		Description: "New description",
		Confidence:  "high",
		Source:      "conversation",
		CreatedAt:   "2026-05-24T10:00:00Z",
		Tags:        []string{"backend"},
		Company:     "New Corp",
		Period:      "2023-01 to 2024-01",
		SourceDetail: "From conversation",
	}

	patch := &extraction.Patch{
		PatchID: "patch-test-002",
		Kind:    "fact_upsert",
		Status:  "proposed",
		Operations: []extraction.Operation{
			{
				OpID:    "op-001",
				Type:    "upsert_fact",
				Target:  "facts/experiences.yaml",
				FactID:  newFact.FactID,
				NewFact: newFact,
			},
		},
	}

	result, err := applier.ApplyPatch(patch)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}

	if !result.Success() {
		t.Fatalf("Patch apply failed: %v", result.FailedOps)
	}

	// Verify file now has both facts
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

	patch := &extraction.Patch{
		Status: "invalid_status",
	}

	_, err := applier.ApplyPatch(patch)
	if err == nil {
		t.Error("ApplyPatch should fail for invalid status")
	}
}

func TestApplyPatch_NoOperations(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "careervault-test-")
	defer os.RemoveAll(tempDir)

	applier := NewApplier(tempDir)

	patch := &extraction.Patch{
		Status:     "proposed",
		Operations: []extraction.Operation{},
	}

	_, err := applier.ApplyPatch(patch)
	if err == nil {
		t.Error("ApplyPatch should fail for patch with no operations")
	}
}
