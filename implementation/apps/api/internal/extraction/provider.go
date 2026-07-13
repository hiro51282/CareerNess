package extraction

import (
	"context"
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
// workflow: conversation → provider → validation → normalization
//
// patch 生成は本パイプラインの責務ではない。呼び出し側（handler 層）が
// YAMLFacts を patch.BuildFactUpsert に渡して patch proposal を組み立てる。
// これにより extraction → patch の import 循環を避ける。
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

	// 抽出 0 件は正常（非 fact の発言では reply のみで会話を継続する）。
	// 0 件時の reply 必須は ValidateAPIResult が保証済み。
	return &ExtractionPipelineResult{
		Reply:             apiResult.Reply,
		ExtractionQuality: apiResult.ExtractionQuality,
		YAMLFacts:         yamlFacts,
	}, nil
}
