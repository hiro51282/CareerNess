package extraction

import (
	"testing"
)

func TestValidateAPIResult(t *testing.T) {
	tests := []struct {
		name    string
		result  *ExtractedFactResult
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil result",
			result:  nil,
			wantErr: true,
			errMsg:  "nil result",
		},
		{
			name: "no facts extracted",
			result: &ExtractedFactResult{
				ExtractedFacts: []ExtractedFact{},
				ExtractionQuality: ExtractionQuality{
					OverallConfidence: "high",
				},
			},
			wantErr: true,
			errMsg:  "no facts extracted",
		},
		{
			name: "valid single fact",
			result: &ExtractedFactResult{
				ExtractedFacts: []ExtractedFact{
					{
						Type:        "experience",
						FactIDHint:  "payment-platform",
						Summary:     "Payment platform redesign",
						Description: "Led redesign of payment system",
						Confidence:  "high",
						Period: &PeriodInfo{
							Start: "2022-04",
							End:   "2023-01",
						},
						ExtractionNotes: []string{
							"User mentioned payment platform work",
						},
						ClarificationQuestions: []string{},
					},
				},
				ExtractionQuality: ExtractionQuality{
					OverallConfidence:       "high",
					Completeness:            "medium",
					NeedsClarificationCount: 0,
					Summary:                 "Extracted payment platform experience",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid fact type",
			result: &ExtractedFactResult{
				ExtractedFacts: []ExtractedFact{
					{
						Type:        "invalid_type",
						FactIDHint:  "test",
						Summary:     "Test",
						Description: "Test description",
						Confidence:  "high",
					},
				},
				ExtractionQuality: ExtractionQuality{
					OverallConfidence: "high",
				},
			},
			wantErr: true,
			errMsg:  "invalid type",
		},
		{
			name: "missing summary",
			result: &ExtractedFactResult{
				ExtractedFacts: []ExtractedFact{
					{
						Type:        "experience",
						FactIDHint:  "test",
						Summary:     "",
						Description: "Test description",
						Confidence:  "high",
					},
				},
				ExtractionQuality: ExtractionQuality{
					OverallConfidence: "high",
				},
			},
			wantErr: true,
			errMsg:  "missing summary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIResult(tt.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateAPIResult() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateYAMLFact(t *testing.T) {
	tests := []struct {
		name    string
		fact    *YAMLFact
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid experience fact",
			fact: &YAMLFact{
				FactID:      "fact-proj-payment-platform",
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Payment platform redesign",
				Description: "Led redesign of payment system",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend", "platform"},
				Company:     "Example Corp",
				Period:      "2022-04 to 2023-01",
				SourceDetail: "User mentioned this project",
			},
			wantErr: false,
		},
		{
			name: "valid achievement fact",
			fact: &YAMLFact{
				FactID:       "fact-ach-latency-improvement",
				Type:         "achievement",
				Status:       "proposed",
				Summary:      "Reduced latency by 30%",
				Description:  "Optimized payment service latency",
				Confidence:   "high",
				Source:       "conversation",
				CreatedAt:    "2026-05-24T10:00:00Z",
				Tags:         []string{"performance"},
				SourceDetail: "From conversation on 2026-05-24",
			},
			wantErr: false,
		},
		{
			name: "valid skill fact",
			fact: &YAMLFact{
				FactID:       "fact-skill-go",
				Type:         "skill",
				Status:       "proposed",
				Summary:      "Go programming",
				Description:  "10 years Go experience",
				Confidence:   "high",
				Source:       "manual",
				CreatedAt:    "2026-05-24T10:00:00Z",
				Tags:         []string{"language"},
				SourceDetail: "Manually added by user",
			},
			wantErr: false,
		},
		{
			name: "missing fact_id",
			fact: &YAMLFact{
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
			},
			wantErr: true,
			errMsg:  "missing fact_id",
		},
		{
			name: "invalid type",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "invalid",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
			},
			wantErr: true,
			errMsg:  "invalid type",
		},
		{
			name: "invalid status",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "experience",
				Status:      "invalid",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
				Company:     "Corp",
				Period:      "2022-01 to 2023-01",
			},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name: "extracted fact with confirmed status (invalid)",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "experience",
				Status:      "confirmed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
				Company:     "Corp",
				Period:      "2022-01 to 2023-01",
			},
			wantErr: true,
			errMsg:  "source=conversation) must have status=proposed",
		},
		{
			name: "experience missing company",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
				Company:     "", // explicitly empty
				Period:      "2022-01 to 2023-01",
			},
			wantErr: true,
			errMsg:  "must have company",
		},
		{
			name: "experience missing period",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
				Company:     "Corp",
				Period:      "", // explicitly empty
			},
			wantErr: true,
			errMsg:  "must have period",
		},
		{
			name: "invalid fact_id format",
			fact: &YAMLFact{
				FactID:      "invalid-format",
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{"backend"},
				Company:     "Corp",
				Period:      "2022-01 to 2023-01",
			},
			wantErr: true,
			errMsg:  "invalid fact_id format",
		},
		{
			name: "missing tags",
			fact: &YAMLFact{
				FactID:      "fact-proj-test",
				Type:        "experience",
				Status:      "proposed",
				Summary:     "Test",
				Description: "Test",
				Confidence:  "high",
				Source:      "conversation",
				CreatedAt:   "2026-05-24T10:00:00Z",
				Tags:        []string{},
				Company:     "Corp",
				Period:      "2022-01 to 2023-01",
			},
			wantErr: true,
			errMsg:  "at least one tag required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateYAMLFact(tt.fact)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateYAMLFact() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateYAMLFact() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func contains(s, substr string) bool {
	// Simple substring check
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
