package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"careerness/api/internal/extraction"
	"gopkg.in/yaml.v3"
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
func (a *Applier) ApplyPatch(p *Patch) (*ApplyResult, error) {
	if p == nil {
		return nil, fmt.Errorf("nil patch")
	}

	if p.Status != StatusProposed && p.Status != StatusApproved {
		return nil, fmt.Errorf("cannot apply patch with status %q (expected proposed or approved)", p.Status)
	}

	if len(p.Operations) == 0 {
		return nil, fmt.Errorf("patch has no operations")
	}

	// Process each operation
	result := &ApplyResult{
		PatchID:      p.PatchID,
		AppliedCount: 0,
		FailedOps:    []string{},
		UpdatedFacts: []*extraction.YAMLFact{},
		AppliedAt:    time.Now().UTC().Format(time.RFC3339) + "Z",
	}

	for i := range p.Operations {
		op := &p.Operations[i]
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
func (a *Applier) applyOperation(op *Operation, result *ApplyResult) error {
	switch op.Type {
	case OpUpsertFact:
		return a.applyUpsertFact(op, result)

	case OpMarkFactStatus:
		return a.applyMarkFactStatus(op, result)

	default:
		return fmt.Errorf("unsupported operation type: %q", op.Type)
	}
}

// factFromChange は operation の change.after を YAMLFact に復元する。
// after は in-process では *extraction.YAMLFact、JSON 経由では map[string]interface{}
// になるため、JSON 往復で両方を吸収する。
func factFromChange(op *Operation) (*extraction.YAMLFact, error) {
	if op.Change.After == nil {
		return nil, fmt.Errorf("upsert_fact requires change.after")
	}
	raw, err := json.Marshal(op.Change.After)
	if err != nil {
		return nil, fmt.Errorf("marshal change.after: %w", err)
	}
	var fact extraction.YAMLFact
	if err := json.Unmarshal(raw, &fact); err != nil {
		return nil, fmt.Errorf("decode change.after into fact: %w", err)
	}
	if fact.FactID == "" {
		// entity_id を fact_id の正本とする（docs: fact 操作は entity_id を持つ）
		fact.FactID = op.EntityID
	}
	if fact.FactID == "" {
		return nil, fmt.Errorf("upsert_fact requires fact_id (entity_id or change.after.fact_id)")
	}
	return &fact, nil
}

// applyUpsertFact appends or updates a fact in the target file
func (a *Applier) applyUpsertFact(op *Operation, result *ApplyResult) error {
	fact, err := factFromChange(op)
	if err != nil {
		return err
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
		if existing.FactID == fact.FactID {
			// Update existing
			facts[i] = fact
			found = true
			break
		}
	}

	if !found {
		// Append new
		facts = append(facts, fact)
	}

	// Write back to file
	content, err := yaml.Marshal(facts)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	result.UpdatedFacts = append(result.UpdatedFacts, fact)
	return nil
}

// applyMarkFactStatus updates a fact's status in the target file.
// 新しい status は operation の fact_status_after を正本とする
// （docs/implementation/ai/ai-patch-format.md）。
func (a *Applier) applyMarkFactStatus(op *Operation, result *ApplyResult) error {
	factID := op.EntityID
	if factID == "" {
		return fmt.Errorf("mark_fact_status requires entity_id")
	}
	if op.FactStatusAfter == "" {
		return fmt.Errorf("mark_fact_status requires fact_status_after")
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
		if fact.FactID == factID {
			fact.Status = op.FactStatusAfter
			fact.UpdatedAt = time.Now().UTC().Format(time.RFC3339) + "Z"

			facts[i] = fact
			found = true
			result.UpdatedFacts = append(result.UpdatedFacts, fact)
			break
		}
	}

	if !found {
		return fmt.Errorf("fact %q not found in %s", factID, op.Target)
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
	PatchID      string                 `json:"patch_id"`
	AppliedCount int                    `json:"applied_count"`
	FailedOps    []string               `json:"failed_ops"`
	UpdatedFacts []*extraction.YAMLFact `json:"updated_facts"`
	AppliedAt    string                 `json:"applied_at"`
	Error        string                 `json:"error"`
}

// Success returns true if apply succeeded (all operations passed)
func (r *ApplyResult) Success() bool {
	return len(r.FailedOps) == 0
}
