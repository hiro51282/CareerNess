package workspace

import (
	"fmt"
	"strings"

	"careerness/api/internal/extraction"
)

// ValidateFact validates a YAMLFact against MVP schema
// This is the authoritative schema validation for all facts in the system
func ValidateFact(fact *extraction.YAMLFact) error {
	// Required fields check
	if err := validateRequiredFields(fact); err != nil {
		return err
	}

	// Type-specific validation
	if err := validateTypeSpecific(fact); err != nil {
		return err
	}

	// Format validation
	if err := validateFormats(fact); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields checks all required fields exist and are non-empty
func validateRequiredFields(fact *extraction.YAMLFact) error {
	requiredFields := map[string]interface{}{
		"fact_id":   fact.FactID,
		"type":      fact.Type,
		"status":    fact.Status,
		"summary":   fact.Summary,
		"description": fact.Description,
		"confidence": fact.Confidence,
		"source":    fact.Source,
		"created_at": fact.CreatedAt,
	}

	for fieldName, fieldValue := range requiredFields {
		if fieldValue == "" {
			return fmt.Errorf("missing or empty required field: %s", fieldName)
		}
	}

	// Tags must have at least one element
	if len(fact.Tags) == 0 {
		return fmt.Errorf("missing required field: tags (at least one tag required)")
	}

	return nil
}

// validateTypeSpecific performs type-specific validation
func validateTypeSpecific(fact *extraction.YAMLFact) error {
	// Validate type value
	validTypes := []string{"experience", "achievement", "skill"}
	if !stringInList(fact.Type, validTypes) {
		return fmt.Errorf("invalid type: %q (must be one of: experience, achievement, skill)", fact.Type)
	}

	// Validate status value
	validStatuses := []string{"confirmed", "proposed"}
	if !stringInList(fact.Status, validStatuses) {
		return fmt.Errorf("invalid status: %q (must be confirmed or proposed)", fact.Status)
	}

	// Validate confidence value
	validConfidence := []string{"high", "medium", "low"}
	if !stringInList(fact.Confidence, validConfidence) {
		return fmt.Errorf("invalid confidence: %q (must be high, medium, or low)", fact.Confidence)
	}

	// Validate source value
	validSources := []string{"conversation", "manual", "import"}
	if !stringInList(fact.Source, validSources) {
		return fmt.Errorf("invalid source: %q (must be one of: conversation, manual, import)", fact.Source)
	}

	// Extracted facts (source=conversation) must always be proposed
	if fact.Source == "conversation" && fact.Status != "proposed" {
		return fmt.Errorf(
			"extracted fact (source=conversation) must have status=proposed, got %q",
			fact.Status,
		)
	}

	// Type-specific required fields
	switch fact.Type {
	case "experience":
		if fact.Company == "" {
			return fmt.Errorf("experience fact must have non-empty company field")
		}
		if fact.Period == "" {
			return fmt.Errorf("experience fact must have non-empty period field")
		}

	case "achievement":
		// Achievement only requires common fields
		// No additional required fields

	case "skill":
		// Skill only requires common fields
		// No additional required fields
	}

	return nil
}

// validateFormats validates the format of specific fields
func validateFormats(fact *extraction.YAMLFact) error {
	// Validate fact_id format
	// Expected: fact-{proj|ach|skill}-{name}
	if !validFactIDFormat(fact.FactID) {
		return fmt.Errorf(
			"invalid fact_id format: %q\n  expected: fact-{type}-{name}\n  example: fact-proj-payment-platform",
			fact.FactID,
		)
	}

	// Validate created_at is non-empty RFC3339 format
	if !validTimestampFormat(fact.CreatedAt) {
		return fmt.Errorf("invalid created_at format: %q (expected RFC3339, e.g., 2026-05-24T10:00:00Z)", fact.CreatedAt)
	}

	// If present, validate updated_at
	if fact.UpdatedAt != "" && !validTimestampFormat(fact.UpdatedAt) {
		return fmt.Errorf("invalid updated_at format: %q (expected RFC3339, e.g., 2026-05-24T10:00:00Z)", fact.UpdatedAt)
	}

	// Validate period format for experience
	if fact.Type == "experience" && fact.Period != "" {
		if !validPeriodFormat(fact.Period) {
			return fmt.Errorf(
				"invalid period format: %q\n  expected: YYYY-MM to YYYY-MM or unknown to unknown\n  example: 2022-04 to 2023-01",
				fact.Period,
			)
		}
	}

	return nil
}

// validFactIDFormat validates fact_id matches pattern: fact-{type}-{name}
func validFactIDFormat(id string) bool {
	// Minimal validation: should start with fact-, have at least 3 parts
	parts := strings.Split(id, "-")
	if len(parts) < 3 {
		return false
	}
	if parts[0] != "fact" {
		return false
	}

	// Type prefix must be valid
	validTypePrefix := []string{"proj", "ach", "skill"}
	if !stringInList(parts[1], validTypePrefix) {
		return false
	}

	// Name must not be empty
	if len(strings.Join(parts[2:], "-")) == 0 {
		return false
	}

	return true
}

// validTimestampFormat validates RFC3339 timestamp format
func validTimestampFormat(timestamp string) bool {
	// Simple check: should contain T and end with Z or have timezone
	// More rigorous parsing could be done with time.Parse
	return strings.Contains(timestamp, "T") && (strings.HasSuffix(timestamp, "Z") || strings.Contains(timestamp, "+"))
}

// validPeriodFormat validates period format: YYYY-MM to YYYY-MM
func validPeriodFormat(period string) bool {
	// Expected: "YYYY-MM to YYYY-MM" or "unknown to unknown"
	parts := strings.Split(period, " to ")
	if len(parts) != 2 {
		return false
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "unknown" {
			continue
		}
		// Check YYYY-MM format (simple check)
		if len(part) != 7 || part[4] != '-' {
			return false
		}
	}

	return true
}

// Helper function
func stringInList(s string, list []string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}
