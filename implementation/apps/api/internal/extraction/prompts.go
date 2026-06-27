package extraction

import "fmt"

// このファイルの prompt は docs/implementation/ai/extraction-specification.md §5 を
// 唯一の正本（SSOT）として逐語転記したもの。独自の prompt engineering は行わない。
// 文言を変更する場合は必ず仕様書を先に更新すること。

// getSystemPrompt returns the system prompt for fact extraction.
// SSOT: extraction-specification.md §5 getSystemPrompt
func getSystemPrompt() string {
	return `You are a career fact extraction agent for CareerNess.

Your task:
- Extract structured career facts from user conversation
- Be conservative: only extract what user explicitly stated
- Mark confidence (high | medium | low) for each field
- Identify uncertainty explicitly
- Generate clarification questions for incomplete information

Do NOT:
- Invent or infer facts not stated
- Mark inferred information as confirmed
- Add metrics the user didn't state
- Fill in missing fields with assumptions

Output ONLY valid JSON matching the provided schema. No markdown, no explanation.`
}

// getUserPrompt returns the user prompt with conversation embedded.
// SSOT: extraction-specification.md §5 getUserPrompt
func getUserPrompt(conversation string) string {
	return fmt.Sprintf(`Extract career facts from this user statement.

User statement:
"%s"

Output JSON with structure:
{
  "extracted_facts": [
    {
      "type": "experience | achievement | skill",
      "fact_id_hint": "semantic-id",
      "summary": "One-liner summary",
      "period": {"start": "YYYY-MM or unknown", "end": "YYYY-MM or unknown"},
      "company": "Company name or null",
      "description": "Detailed description",
      "confidence": "high | medium | low",
      "details": {
        // Type-specific details
      },
      "extraction_notes": ["Explicit facts"],
      "clarification_questions": ["Questions to ask"]
    }
  ],
  "extraction_quality": {
    "overall_confidence": "high | medium | low",
    "completeness": "high | medium | low",
    "needs_clarification_count": <number>,
    "summary": "Brief summary"
  }
}`, conversation)
}
