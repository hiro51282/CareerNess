package extraction

import (
	"strings"
	"testing"
	"time"
)

// TestNormalizeToYAMLFact_CreatedAtIsValidRFC3339 は created_at が妥当な RFC3339 で
// 生成されること（以前の "...ZZ" 不正形式が再発しないこと）を検証する。
func TestNormalizeToYAMLFact_CreatedAtIsValidRFC3339(t *testing.T) {
	extracted := &ExtractedFact{
		Type:        "skill",
		FactIDHint:  "go",
		Summary:     "Go",
		Description: "Go programming",
		Confidence:  "high",
	}

	fact, err := NormalizeToYAMLFact(extracted)
	if err != nil {
		t.Fatalf("NormalizeToYAMLFact() error = %v", err)
	}

	if strings.HasSuffix(fact.CreatedAt, "ZZ") {
		t.Errorf("created_at が不正な二重 Z 形式: %q", fact.CreatedAt)
	}
	if _, err := time.Parse(time.RFC3339, fact.CreatedAt); err != nil {
		t.Errorf("created_at が RFC3339 として解析できない: %q (%v)", fact.CreatedAt, err)
	}
}

func TestNormalizeToYAMLFact_Experience(t *testing.T) {
	extracted := &ExtractedFact{
		Type:        "experience",
		FactIDHint:  "payment-platform",
		Summary:     "Payment platform redesign",
		Company:     "Example Corp",
		Description: "Led redesign of payment system from Python to Go",
		Confidence:  "high",
		Period: &PeriodInfo{
			Start: "2022-04",
			End:   "2023-01",
		},
		Details: map[string]interface{}{
			"tech_stack": []interface{}{"Go", "gRPC", "PostgreSQL"},
			"team_size":  float64(8),
			"impact": map[string]interface{}{
				"latency_improvement_pct": float64(30),
			},
		},
		ExtractionNotes: []string{
			"User mentioned payment platform redesign",
			"Timeline confirmed as 2022-04 to 2023-01",
		},
		ClarificationQuestions: []string{
			"What was the baseline latency before improvement?",
		},
	}

	fact, err := NormalizeToYAMLFact(extracted)
	if err != nil {
		t.Fatalf("NormalizeToYAMLFact() error = %v", err)
	}

	// Verify basic fields
	if fact.FactID != "fact-proj-payment-platform" {
		t.Errorf("FactID = %q, want fact-proj-payment-platform", fact.FactID)
	}
	if fact.Type != "experience" {
		t.Errorf("Type = %q, want experience", fact.Type)
	}
	if fact.Status != "proposed" {
		t.Errorf("Status = %q, want proposed", fact.Status)
	}
	if fact.Company != "Example Corp" {
		t.Errorf("Company = %q, want Example Corp", fact.Company)
	}
	if fact.Period != "2022-04 to 2023-01" {
		t.Errorf("Period = %q, want 2022-04 to 2023-01", fact.Period)
	}

	// Verify tech_stack
	if len(fact.TechStack) != 3 {
		t.Errorf("TechStack length = %d, want 3", len(fact.TechStack))
	}

	// Verify team_size
	if fact.TeamSize == nil || *fact.TeamSize != 8 {
		t.Errorf("TeamSize = %v, want 8", fact.TeamSize)
	}

	// Verify impact
	if fact.Impact == nil {
		t.Error("Impact should not be nil")
	}

	// Verify source_detail contains clarification questions
	if !strings.Contains(fact.SourceDetail, "Clarification needed") {
		t.Error("SourceDetail should contain clarification questions")
	}
}

func TestNormalizeToYAMLFact_Achievement(t *testing.T) {
	extracted := &ExtractedFact{
		Type:        "achievement",
		FactIDHint:  "latency-improvement",
		Summary:     "Reduced latency by 30%",
		Description: "Optimized payment service latency from 200ms to 140ms",
		Confidence:  "high",
		Details: map[string]interface{}{
			"impact": map[string]interface{}{
				"latency_improvement_pct": float64(30),
				"user_experience":         "faster payments",
			},
		},
		ExtractionNotes: []string{
			"Latency reduction confirmed from prod logs",
		},
		ClarificationQuestions: []string{},
	}

	fact, err := NormalizeToYAMLFact(extracted)
	if err != nil {
		t.Fatalf("NormalizeToYAMLFact() error = %v", err)
	}

	if fact.FactID != "fact-ach-latency-improvement" {
		t.Errorf("FactID = %q, want fact-ach-latency-improvement", fact.FactID)
	}
	if fact.Type != "achievement" {
		t.Errorf("Type = %q, want achievement", fact.Type)
	}
	if fact.Status != "proposed" {
		t.Errorf("Status = %q, want proposed", fact.Status)
	}

	// Verify impact
	if fact.Impact == nil {
		t.Error("Impact should not be nil")
	}
}

func TestNormalizeToYAMLFact_Skill(t *testing.T) {
	extracted := &ExtractedFact{
		Type:        "skill",
		FactIDHint:  "go-systems",
		Summary:     "Go systems programming",
		Description: "10+ years experience with Go, expertise in concurrent systems",
		Confidence:  "high",
		ExtractionNotes: []string{
			"User stated 10+ years Go experience",
		},
		ClarificationQuestions: []string{},
	}

	fact, err := NormalizeToYAMLFact(extracted)
	if err != nil {
		t.Fatalf("NormalizeToYAMLFact() error = %v", err)
	}

	if fact.FactID != "fact-skill-go-systems" {
		t.Errorf("FactID = %q, want fact-skill-go-systems", fact.FactID)
	}
	if fact.Type != "skill" {
		t.Errorf("Type = %q, want skill", fact.Type)
	}
}

func TestNormalizeToYAMLFact_DefaultCompany(t *testing.T) {
	extracted := &ExtractedFact{
		Type:        "experience",
		FactIDHint:  "test",
		Summary:     "Test project",
		Company:     "", // Empty company
		Description: "Test description",
		Confidence:  "high",
		Period: &PeriodInfo{
			Start: "2022-01",
			End:   "2023-01",
		},
		ExtractionNotes:        []string{},
		ClarificationQuestions: []string{},
	}

	fact, err := NormalizeToYAMLFact(extracted)
	if err != nil {
		t.Fatalf("NormalizeToYAMLFact() error = %v", err)
	}

	if fact.Company != "未確認" {
		t.Errorf("Company = %q, want 未確認 (default for missing)", fact.Company)
	}
}

func TestGenerateFactID(t *testing.T) {
	tests := []struct {
		name     string
		fact     *ExtractedFact
		expected string
	}{
		{
			name: "experience",
			fact: &ExtractedFact{
				Type:       "experience",
				FactIDHint: "payment-platform",
			},
			expected: "fact-proj-payment-platform",
		},
		{
			name: "achievement",
			fact: &ExtractedFact{
				Type:       "achievement",
				FactIDHint: "latency-improvement",
			},
			expected: "fact-ach-latency-improvement",
		},
		{
			name: "skill",
			fact: &ExtractedFact{
				Type:       "skill",
				FactIDHint: "go-systems",
			},
			expected: "fact-skill-go-systems",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFactID(tt.fact)
			if result != tt.expected {
				t.Errorf("generateFactID() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestInferTags(t *testing.T) {
	tests := []struct {
		name      string
		fact      *ExtractedFact
		wantTags  []string
		wantCount int
	}{
		{
			name: "backend keywords",
			fact: &ExtractedFact{
				Summary:     "Backend API development",
				Description: "Developed payment API and database schema",
				Type:        "experience",
			},
			wantTags:  []string{"backend"},
			wantCount: 1,
		},
		{
			name: "platform and performance keywords",
			fact: &ExtractedFact{
				Summary:     "Infrastructure optimization",
				Description: "Kubernetes deployment and performance tuning",
				Type:        "achievement",
			},
			wantTags:  []string{"performance", "platform"},
			wantCount: 2,
		},
		{
			name: "leadership keywords",
			fact: &ExtractedFact{
				Summary:     "Team leadership",
				Description: "Led technical coordination and mentoring",
				Type:        "achievement",
			},
			wantTags:  []string{"leadership"},
			wantCount: 1,
		},
		{
			name: "multiple tags",
			fact: &ExtractedFact{
				Summary:     "Go backend optimization",
				Description: "Kubernetes platform with performance improvements and team coordination",
				Type:        "experience",
			},
			wantTags:  []string{"backend", "leadership", "performance", "platform"},
			wantCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := inferTags(tt.fact)

			if len(tags) != tt.wantCount {
				t.Errorf("inferTags() returned %d tags, want %d. Tags: %v", len(tags), tt.wantCount, tags)
			}

			// Verify expected tags are present
			for _, expected := range tt.wantTags {
				found := false
				for _, tag := range tags {
					if tag == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("inferTags() missing expected tag %q. Got: %v", expected, tags)
				}
			}
		})
	}
}

func TestFormatSourceDetail(t *testing.T) {
	extracted := &ExtractedFact{
		ExtractionNotes: []string{
			"User mentioned payment platform work",
			"Timeline 2022-04 to 2023-01 confirmed",
		},
		ClarificationQuestions: []string{
			"What was the exact team size?",
			"Can you confirm the latency baseline?",
		},
	}

	result := formatSourceDetail(extracted)

	// Verify structure
	if !strings.Contains(result, "Extraction notes:") {
		t.Error("formatSourceDetail should contain 'Extraction notes:'")
	}
	if !strings.Contains(result, "Clarification needed:") {
		t.Error("formatSourceDetail should contain 'Clarification needed:'")
	}

	// Verify content
	for _, note := range extracted.ExtractionNotes {
		if !strings.Contains(result, note) {
			t.Errorf("formatSourceDetail should contain extraction note: %q", note)
		}
	}

	for _, q := range extracted.ClarificationQuestions {
		if !strings.Contains(result, q) {
			t.Errorf("formatSourceDetail should contain clarification question: %q", q)
		}
	}
}

func TestFormatPeriod(t *testing.T) {
	tests := []struct {
		name     string
		period   *PeriodInfo
		expected string
	}{
		{
			name: "valid period",
			period: &PeriodInfo{
				Start: "2022-04",
				End:   "2023-01",
			},
			expected: "2022-04 to 2023-01",
		},
		{
			name: "unknown start",
			period: &PeriodInfo{
				Start: "unknown",
				End:   "2023-01",
			},
			expected: "unknown to 2023-01",
		},
		{
			name:     "nil period",
			period:   nil,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPeriod(tt.period)
			if result != tt.expected {
				t.Errorf("formatPeriod() = %q, want %q", result, tt.expected)
			}
		})
	}
}
