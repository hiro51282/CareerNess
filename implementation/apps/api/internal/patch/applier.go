package patch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"careerness/api/internal/extraction"
)

// Applier は patch を workspace に apply する
type Applier struct {
	workspacePath string
}

// NewApplier creates a new patch applier
func NewApplier(workspacePath string) *Applier {
	return &Applier{
		workspacePath: workspacePath,
	}
}

// ApplyPatch applies a patch to the workspace
// Returns updated facts and an error if apply fails
func (a *Applier) ApplyPatch(patch *extraction.Patch) (*ApplyResult, error) {
	if patch == nil {
		return nil, fmt.Errorf("nil patch")
	}

	if patch.Status != "proposed" && patch.Status != "approved" {
		return nil, fmt.Errorf("cannot apply patch with status %q (expected proposed or approved)", patch.Status)
	}

	if len(patch.Operations) == 0 {
		return nil, fmt.Errorf("patch has no operations")
	}

	// Process each operation
	result := &ApplyResult{
		PatchID:      patch.PatchID,
		AppliedCount: 0,
		FailedOps:    []string{},
		UpdatedFacts: []*extraction.YAMLFact{},
		AppliedAt:    time.Now().UTC().Format(time.RFC3339) + "Z",
	}

	for i := range patch.Operations {
		op := &patch.Operations[i]
		if err := a.applyOperation(op, result); err != nil {
			result.FailedOps = append(result.FailedOps, fmt.Sprintf("%s: %v", op.OpID, err))
		} else {
			result.AppliedCount++
		}
	}

	// Check if any operations succeeded
	if result.AppliedCount == 0 {
		return result, fmt.Errorf("all operations failed")
	}

	return result, nil
}

// applyOperation applies a single operation
func (a *Applier) applyOperation(op *extraction.Operation, result *ApplyResult) error {
	switch op.Type {
	case "upsert_fact":
		return a.applyUpsertFact(op, result)

	case "update_fact_status":
		return a.applyUpdateFactStatus(op, result)

	default:
		return fmt.Errorf("unknown operation type: %q", op.Type)
	}
}

// applyUpsertFact appends or updates a fact in the target file
func (a *Applier) applyUpsertFact(op *extraction.Operation, result *ApplyResult) error {
	if op.NewFact == nil {
		return fmt.Errorf("upsert_fact requires new_fact")
	}

	targetPath := filepath.Join(a.workspacePath, op.Target)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Read existing facts (if file exists)
	facts := []*extraction.YAMLFact{}
	if _, err := os.Stat(targetPath); err == nil {
		// File exists, read it
		content, err := os.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		if err := yaml.Unmarshal(content, &facts); err != nil {
			return fmt.Errorf("parse yaml: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat file: %w", err)
	}
	// If file doesn't exist, facts will be empty slice

	// Check if fact already exists (by ID)
	found := false
	for i, existing := range facts {
		if existing.FactID == op.NewFact.FactID {
			// Update existing
			facts[i] = op.NewFact
			found = true
			break
		}
	}

	if !found {
		// Append new
		facts = append(facts, op.NewFact)
	}

	// Write back to file
	content, err := yaml.Marshal(facts)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	result.UpdatedFacts = append(result.UpdatedFacts, op.NewFact)
	return nil
}

// applyUpdateFactStatus updates a fact's status in the target file
func (a *Applier) applyUpdateFactStatus(op *extraction.Operation, result *ApplyResult) error {
	if op.FactID == "" {
		return fmt.Errorf("update_fact_status requires fact_id")
	}

	targetPath := filepath.Join(a.workspacePath, op.Target)

	// Read existing facts
	content, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	facts := []*extraction.YAMLFact{}
	if err := yaml.Unmarshal(content, &facts); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	// Find and update fact
	found := false
	for i, fact := range facts {
		if fact.FactID == op.FactID {
			// Extract new status from operation (it should be in change.after or similar)
			// For MVP, we'll look in the operation details
			// This is simplified; production version would parse operation.Change

			fact.Status = "confirmed" // Simplified: always set to confirmed
			fact.UpdatedAt = time.Now().UTC().Format(time.RFC3339) + "Z"

			facts[i] = fact
			found = true
			result.UpdatedFacts = append(result.UpdatedFacts, fact)
			break
		}
	}

	if !found {
		return fmt.Errorf("fact %q not found in %s", op.FactID, op.Target)
	}

	// Write back
	newContent, err := yaml.Marshal(facts)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(targetPath, newContent, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// ApplyResult は patch apply の結果
type ApplyResult struct {
	PatchID      string                   `json:"patch_id"`
	AppliedCount int                      `json:"applied_count"`
	FailedOps    []string                 `json:"failed_ops"`
	UpdatedFacts []*extraction.YAMLFact   `json:"updated_facts"`
	AppliedAt    string                   `json:"applied_at"`
	Error        string                   `json:"error"`
}

// Success returns true if apply succeeded (all operations passed)
func (r *ApplyResult) Success() bool {
	return len(r.FailedOps) == 0
}
