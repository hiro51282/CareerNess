package extraction

import (
	"context"
)

// MockExtractionProvider is a testing/MVP implementation that returns mock results
// In production, this would be replaced with Codex app-server integration
type MockExtractionProvider struct {
	ResponseOverride *ExtractedFactResult // For testing
}

// NewMockExtractionProvider creates a new mock provider
func NewMockExtractionProvider() *MockExtractionProvider {
	return &MockExtractionProvider{}
}

// ExtractFacts returns pre-canned mock results
// In production, this would call Codex app-server /extract endpoint
func (m *MockExtractionProvider) ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error) {
	// If override is set (for testing), return it
	if m.ResponseOverride != nil {
		return m.ResponseOverride, nil
	}

	// Default mock response based on conversation keywords
	// This is sufficient for MVP testing and demonstration

	// Simple keyword-based extraction (not ML-powered, obviously)
	// In production: POST to Codex app-server, which calls the user's LLM provider

	if len(conversation) == 0 {
		return &ExtractedFactResult{
			ExtractedFacts: []ExtractedFact{},
			ExtractionQuality: ExtractionQuality{
				OverallConfidence:       "low",
				Completeness:            "low",
				NeedsClarificationCount: 0,
				Summary:                 "Empty conversation",
			},
		}, nil
	}

	// Mock: always return a single experience fact
	fact := ExtractedFact{
		Type:        "experience",
		FactIDHint:  "mock-project",
		Summary:     "Mock project from conversation",
		Company:     "Mock Corp",
		Description: "This is a mock extraction result for testing",
		Confidence:  "low",
		Period: &PeriodInfo{
			Start: "2026-01-01",
			End:   "2026-05-24",
		},
		Details: map[string]interface{}{
			"tech_stack": []interface{}{"Go", "YAML"},
		},
		ExtractionNotes: []string{
			"This is a mock result - actual extraction would come from LLM",
		},
		ClarificationQuestions: []string{
			"Is this a real project or a test?",
		},
	}

	return &ExtractedFactResult{
		ExtractedFacts: []ExtractedFact{fact},
		ExtractionQuality: ExtractionQuality{
			OverallConfidence:       "low",
			Completeness:            "low",
			NeedsClarificationCount: 1,
			Summary:                 "Mock extraction - actual LLM call would come from Codex app-server",
		},
	}, nil
}

// Name returns the provider name
func (m *MockExtractionProvider) Name() string {
	return "mock"
}

// SetResponse allows test to override the response
func (m *MockExtractionProvider) SetResponse(result *ExtractedFactResult) {
	m.ResponseOverride = result
}

// 実 provider（CodexExtractionProvider）は codex_provider.go に実装済み。
// provider 選択（どの provider を返すか）は provider_factory.go に集約した。
// 旧来あった空の reference stub と silent fallback 型の NewProvider は削除した。
