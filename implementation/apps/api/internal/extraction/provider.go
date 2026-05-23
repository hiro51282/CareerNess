package extraction

import (
	"context"
	"fmt"
	"time"
)

// ExtractionProvider は会話からの fact 抽出を担当する provider の interface
// 実装は Codex app-server、OpenAI、Anthropic、またはモックなど、どのプロバイダーでも可能
type ExtractionProvider interface {
	// ExtractFacts calls the LLM provider to extract facts from conversation
	// Returns structured JSON result ready for normalization
	ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error)

	// Provider name for logging/debugging
	Name() string
}

// ExtractionService は extraction workflow を orchestrate する
// Provider に依存せず、structure extraction と validation に集中
type ExtractionService struct {
	provider ExtractionProvider
}

// NewExtractionService creates a new extraction service with a given provider
func NewExtractionService(provider ExtractionProvider) *ExtractionService {
	return &ExtractionService{
		provider: provider,
	}
}

// ExtractFromConversation orchestrates the complete extraction pipeline
// workflow: conversation → provider → validation → normalization → patch generation
func (s *ExtractionService) ExtractFromConversation(
	ctx context.Context,
	conversation string,
	sessionID string,
) (*ExtractionPipelineResult, error) {

	// Step 1: Call provider (LLM abstraction - could be Codex, OpenAI, Anthropic, etc.)
	apiResult, err := s.provider.ExtractFacts(ctx, conversation)
	if err != nil {
		return nil, err
	}

	// Step 2: Validate JSON structure (provider-independent)
	if err := ValidateAPIResult(apiResult); err != nil {
		return nil, err
	}

	// Step 3: Convert to YAML (provider-independent)
	yamlFacts := make([]*YAMLFact, 0, len(apiResult.ExtractedFacts))
	for _, extracted := range apiResult.ExtractedFacts {
		yamlFact, err := NormalizeToYAMLFact(&extracted)
		if err != nil {
			continue
		}

		// Step 4: Validate against MVP schema (provider-independent)
		if err := ValidateYAMLFact(yamlFact); err != nil {
			return nil, err
		}

		yamlFacts = append(yamlFacts, yamlFact)
	}

	if len(yamlFacts) == 0 {
		return nil, err
	}

	// Step 5: Generate patches (provider-independent)
	patches, err := generatePatches(yamlFacts, sessionID)
	if err != nil {
		return nil, err
	}

	return &ExtractionPipelineResult{
		ExtractionQuality: apiResult.ExtractionQuality,
		YAMLFacts:         yamlFacts,
		Patches:           patches,
	}, nil
}

// generatePatches creates patch proposals from YAMLFacts
func generatePatches(facts []*YAMLFact, sessionID string) ([]*Patch, error) {
	patches := make([]*Patch, 0, len(facts))

	for i, fact := range facts {
		patch := generatePatchFromFact(fact, sessionID, i)
		patches = append(patches, patch)
	}

	return patches, nil
}

// generatePatchFromFact creates a patch proposal from a YAMLFact
func generatePatchFromFact(fact *YAMLFact, sessionID string, index int) *Patch {
	now := time.Now().UTC().Format(time.RFC3339) + "Z"
	targetFile := fmt.Sprintf("facts/%ss.yaml", fact.Type)

	patchID := fmt.Sprintf("patch-%s-%03d",
		now[0:10],
		index+1,
	)

	patch := &Patch{
		PatchID:     patchID,
		Kind:        "fact_upsert",
		Status:      "proposed",
		CreatedAt:   now,
		CreatedBy:   "ai",
		WorkspaceID: "local-careervault",
		SessionID:   sessionID,
		Summary: fmt.Sprintf(
			"Extracted %s fact: %s\nConfidence: %s. Status: proposed.",
			fact.Type,
			fact.Summary,
			fact.Confidence,
		),
		Operations: []Operation{
			{
				OpID:    "op-001",
				Type:    "upsert_fact",
				Target:  targetFile,
				FactID:  fact.FactID,
				NewFact: fact,
			},
		},
	}

	return patch
}
