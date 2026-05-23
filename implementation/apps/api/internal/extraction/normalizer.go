package extraction

import (
	"fmt"
	"strings"
	"time"
)

// YAMLFact is defined in models.go

// NormalizeToYAMLFact converts ExtractedFact to YAMLFact (MVP schema format)
func NormalizeToYAMLFact(extracted *ExtractedFact) (*YAMLFact, error) {
	now := time.Now().UTC().Format(time.RFC3339) + "Z"

	// Base fact (common to all types)
	fact := &YAMLFact{
		FactID:       generateFactID(extracted),
		Type:         extracted.Type,
		Status:       "proposed", // Always proposed from extraction
		Summary:      extracted.Summary,
		Description:  extracted.Description,
		Confidence:   extracted.Confidence,
		Source:       "conversation",
		CreatedAt:    now,
		Tags:         inferTags(extracted),
		SourceDetail: formatSourceDetail(extracted),
	}

	// Type-specific normalization
	switch extracted.Type {
	case "experience":
		if extracted.Company != "" {
			fact.Company = extracted.Company
		} else {
			fact.Company = "未確認"
		}

		if extracted.Period != nil {
			fact.Period = formatPeriod(extracted.Period)
		}

		// Extract optional fields from Details
		if extracted.Details != nil {
			if techStack, ok := extracted.Details["tech_stack"].([]interface{}); ok {
				fact.TechStack = interfaceSliceToStrings(techStack)
			}
			if teamSize, ok := extracted.Details["team_size"].(float64); ok {
				size := int(teamSize)
				fact.TeamSize = &size
			}
			if impact, ok := extracted.Details["impact"].(map[string]interface{}); ok {
				fact.Impact = impact
			}
		}

	case "achievement":
		// Extract optional fields from Details
		if extracted.Details != nil {
			if impact, ok := extracted.Details["impact"].(map[string]interface{}); ok {
				fact.Impact = impact
			}
		}

	case "skill":
		// Skill type doesn't require additional processing
		// Details can be stored but are not required
	}

	return fact, nil
}

// generateFactID creates stable fact ID from extracted hint
func generateFactID(fact *ExtractedFact) string {
	switch fact.Type {
	case "experience":
		return fmt.Sprintf("fact-proj-%s", fact.FactIDHint)
	case "achievement":
		return fmt.Sprintf("fact-ach-%s", fact.FactIDHint)
	case "skill":
		return fmt.Sprintf("fact-skill-%s", fact.FactIDHint)
	default:
		return fmt.Sprintf("fact-%s", fact.FactIDHint)
	}
}

// inferTags infers appropriate tags from extracted content
func inferTags(fact *ExtractedFact) []string {
	tags := make(map[string]bool)

	// Combine all text for keyword analysis
	fullText := strings.ToLower(fact.Summary + " " + fact.Description)

	// Keyword-based inference
	keywords := map[string][]string{
		"backend":     {"backend", "api", "server", "database", "go", "python"},
		"frontend":    {"frontend", "ui", "react", "javascript", "typescript", "web"},
		"platform":    {"platform", "kubernetes", "devops", "infra", "docker", "aws"},
		"performance": {"performance", "latency", "optimization", "速度", "最適化"},
		"leadership":  {"lead", "team", "coordination", "mentoring", "リード", "調整"},
		"architecture": {"architecture", "design", "system design", "アーキテクチャ"},
	}

	for tag, kws := range keywords {
		if containsAnyKeyword(fullText, kws) {
			tags[tag] = true
		}
	}

	// Default tags based on type if no keywords matched
	if len(tags) == 0 {
		switch fact.Type {
		case "experience":
			tags["backend"] = true // Conservative default
		case "skill":
			tags["language"] = true
		case "achievement":
			tags["impact-metrics"] = true
		}
	}

	// Convert to sorted slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}

	// Sort for deterministic output
	sortStrings(result)
	return result
}

// formatSourceDetail creates source_detail field from extraction metadata
func formatSourceDetail(fact *ExtractedFact) string {
	var lines []string

	lines = append(lines, "Extraction notes:")
	for _, note := range fact.ExtractionNotes {
		lines = append(lines, fmt.Sprintf("- %s", note))
	}

	if len(fact.ClarificationQuestions) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Clarification needed:")
		for _, q := range fact.ClarificationQuestions {
			lines = append(lines, fmt.Sprintf("- %s", q))
		}
	}

	return strings.Join(lines, "\n")
}

// formatPeriod formats PeriodInfo to string
func formatPeriod(p *PeriodInfo) string {
	if p == nil {
		return "unknown"
	}
	return fmt.Sprintf("%s to %s", p.Start, p.End)
}

// interfaceSliceToStrings converts []interface{} to []string
func interfaceSliceToStrings(slice []interface{}) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// containsAnyKeyword checks if text contains any of the keywords
func containsAnyKeyword(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// sortStrings sorts a slice of strings in place
func sortStrings(s []string) {
	// Simple bubble sort for deterministic output
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}
