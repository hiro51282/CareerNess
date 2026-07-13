package extraction

import (
	"fmt"
	"strings"
)

// ValidateAPIResult validates Claude API output structure
func ValidateAPIResult(result *ExtractedFactResult) error {
	if result == nil {
		return fmt.Errorf("nil result")
	}

	// 抽出 0 件は正常（非 fact の発言）。ただしその場合は会話返信 reply が必須で、
	// facts も reply も無い応答は契約違反とする（extraction-specification.md）。
	if len(result.ExtractedFacts) == 0 && strings.TrimSpace(result.Reply) == "" {
		return fmt.Errorf("empty result: no facts extracted and no reply")
	}

	for i, fact := range result.ExtractedFacts {
		if err := validateExtractedFact(&fact, i); err != nil {
			return err
		}
	}

	// Validate extraction quality
	if result.ExtractionQuality.OverallConfidence == "" {
		return fmt.Errorf("missing extraction_quality.overall_confidence")
	}

	validConfidence := []string{"high", "medium", "low"}
	if !stringInList(result.ExtractionQuality.OverallConfidence, validConfidence) {
		return fmt.Errorf("invalid extraction_quality.overall_confidence: %q", result.ExtractionQuality.OverallConfidence)
	}

	return nil
}

// validateExtractedFact validates a single extracted fact
func validateExtractedFact(fact *ExtractedFact, index int) error {
	// Validate required fields
	if fact.Type == "" {
		return fmt.Errorf("fact[%d]: missing type", index)
	}

	validTypes := []string{"experience", "achievement", "skill"}
	if !stringInList(fact.Type, validTypes) {
		return fmt.Errorf("fact[%d]: invalid type %q", index, fact.Type)
	}

	if fact.FactIDHint == "" {
		return fmt.Errorf("fact[%d]: missing fact_id_hint", index)
	}

	if fact.Summary == "" {
		return fmt.Errorf("fact[%d]: missing summary", index)
	}

	if fact.Description == "" {
		return fmt.Errorf("fact[%d]: missing description", index)
	}

	// Validate confidence
	validConfidence := []string{"high", "medium", "low"}
	if !stringInList(fact.Confidence, validConfidence) {
		return fmt.Errorf("fact[%d]: invalid confidence %q", index, fact.Confidence)
	}

	// Type-specific validation
	if fact.Type == "experience" {
		if fact.Period == nil {
			return fmt.Errorf("fact[%d]: experience must have period", index)
		}
	}

	return nil
}

// ValidateYAMLFact validates a YAMLFact against MVP schema
func ValidateYAMLFact(fact *YAMLFact) error {
	// Required fields
	if fact.FactID == "" {
		return fmt.Errorf("missing fact_id")
	}

	// Validate fact ID format
	if !validFactIDFormat(fact.FactID) {
		return fmt.Errorf("invalid fact_id format: %q (expected fact-{type}-{name})", fact.FactID)
	}

	// Validate type
	validTypes := []string{"experience", "achievement", "skill"}
	if !stringInList(fact.Type, validTypes) {
		return fmt.Errorf("invalid type: %q", fact.Type)
	}

	// Validate status
	validStatuses := []string{"confirmed", "proposed"}
	if !stringInList(fact.Status, validStatuses) {
		return fmt.Errorf("invalid status: %q", fact.Status)
	}

	// Validate confidence
	validConfidence := []string{"high", "medium", "low"}
	if !stringInList(fact.Confidence, validConfidence) {
		return fmt.Errorf("invalid confidence: %q", fact.Confidence)
	}

	// Extracted facts must be proposed
	if fact.Source == "conversation" && fact.Status != "proposed" {
		return fmt.Errorf("extracted fact (source=conversation) must have status=proposed, got %q", fact.Status)
	}

	// Non-empty required fields
	if fact.Summary == "" {
		return fmt.Errorf("missing summary")
	}
	if fact.Description == "" {
		return fmt.Errorf("missing description")
	}
	if fact.Source == "" {
		return fmt.Errorf("missing source")
	}
	if fact.CreatedAt == "" {
		return fmt.Errorf("missing created_at")
	}

	// Type-specific validation
	if fact.Type == "experience" {
		if fact.Company == "" {
			return fmt.Errorf("experience must have company")
		}
		if fact.Period == "" {
			return fmt.Errorf("experience must have period")
		}
	}

	// Validate tags are not empty
	if len(fact.Tags) == 0 {
		return fmt.Errorf("missing tags (at least one tag required)")
	}

	return nil
}

// Helper functions

func stringInList(s string, list []string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}

func validFactIDFormat(id string) bool {
	// Expected format: fact-{type}-{name}
	// where type is one of: proj (experience), ach (achievement), skill (skill)
	parts := strings.Split(id, "-")
	if len(parts) < 3 {
		return false
	}
	if parts[0] != "fact" {
		return false
	}
	validTypePrefix := []string{"proj", "ach", "skill"}
	return stringInList(parts[1], validTypePrefix)
}
