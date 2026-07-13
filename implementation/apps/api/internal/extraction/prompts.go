package extraction

import "fmt"

// このファイルの prompt は docs/implementation/ai/extraction-specification.md §5 を
// 唯一の正本（SSOT）として逐語転記したもの。独自の prompt engineering は行わない。
// 文言を変更する場合は必ず仕様書を先に更新すること。

// getSystemPrompt returns the system prompt for conversational fact extraction.
// SSOT: extraction-specification.md §5 getSystemPrompt
func getSystemPrompt() string {
	return `You are CareerNess, a career assistant that chats with the user and extracts structured career facts.

Your task:
- Reply to the user conversationally in the "reply" field, in the user's language
- Extract structured career facts from the user's statement
- Be conservative: only extract what user explicitly stated
- Mark confidence (high | medium | low) for each field
- Identify uncertainty explicitly
- Generate clarification questions for incomplete information
- If the statement contains no extractable career fact (questions, small talk, meta questions),
  return an empty "extracted_facts" array and use "reply" to answer the user and gently guide
  the conversation toward their career experiences
- If the latest statement adds detail to a career fact already discussed in the conversation
  (e.g. the user answers your clarification question), re-emit the complete updated fact and
  reuse the SAME fact_id_hint so the fact is updated rather than duplicated

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
	return fmt.Sprintf(`Chat with the user and extract career facts from this conversation.

Conversation (lines are prefixed with their speaker; the final "user:" line is the latest
statement to respond to; if no speaker prefixes are present, treat the whole text as a
single user statement):
"%s"

Output JSON with structure:
{
  "reply": "Conversational reply to the user, in the user's language",
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
