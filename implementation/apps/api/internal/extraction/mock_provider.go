package extraction

import (
	"context"
	"strings"
	"unicode/utf8"
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

// ExtractFacts は会話を反映した mock 抽出結果を返す。
// 実 LLM を呼ばない代わりに、発言内容から単一の experience fact を組み立てる
// （既定 provider として、開発・CI が API key 無しで会話 UX を確認できるようにする）。
//
// type の自動分類・action/decision の充足・複数 fact 抽出は実 AI（Codex）の責務であり、
// mock では行わない。confidence=low は出力が mock 由来であることを示す。
func (m *MockExtractionProvider) ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error) {
	// テスト用の上書きがあれば優先する。
	if m.ResponseOverride != nil {
		return m.ResponseOverride, nil
	}

	trimmed := strings.TrimSpace(conversation)
	if trimmed == "" {
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

	// Company は空のままにし、normalizer の既定（"未確認"）に委ねる。
	// fact を捏造しないため tech_stack 等の details は埋めない。
	fact := ExtractedFact{
		Type:        "experience",
		FactIDHint:  mockFactIDHint(trimmed),
		Summary:     mockSummary(trimmed),
		Description: trimmed,
		Confidence:  "low",
		Period:      &PeriodInfo{Start: "unknown", End: "unknown"},
		ExtractionNotes: []string{
			"Mock extraction (not a real LLM call). Set EXTRACTION_PROVIDER=codex for real extraction.",
		},
	}

	return &ExtractedFactResult{
		ExtractedFacts: []ExtractedFact{fact},
		ExtractionQuality: ExtractionQuality{
			OverallConfidence:       "low",
			Completeness:            "low",
			NeedsClarificationCount: 0,
			Summary:                 "Mock extraction from conversation",
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

// mockSummary は会話の先頭を fact summary 用に rune 単位で短く整える。
func mockSummary(conversation string) string {
	s := strings.TrimSpace(conversation)
	const max = 40
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	return string([]rune(s)[:max]) + "…"
}

// mockFactIDHint は会話から fact_id_hint 用の安定 slug を作る。
// normalizer が "fact-proj-<hint>" を生成し validFactIDFormat を満たすよう、
// 英数字以外は "-" に畳み、空になる場合は "project" にフォールバックする。
func mockFactIDHint(conversation string) string {
	runes := []rune(conversation)
	if len(runes) > 12 {
		runes = runes[:12]
	}
	var b strings.Builder
	for _, r := range runes {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r + 32)
		default:
			b.WriteRune('-')
		}
	}
	hint := strings.Trim(b.String(), "-")
	if hint == "" {
		return "project"
	}
	return hint
}

// 実 provider（CodexExtractionProvider）は codex_provider.go に実装済み。
// provider 選択（どの provider を返すか）は provider_factory.go に集約した。
// 旧来あった空の reference stub と silent fallback 型の NewProvider は削除した。
