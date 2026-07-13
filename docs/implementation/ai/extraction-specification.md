# Fact Extraction Pipeline Specification

> **[訂正 2026-05-31]**  
> このドキュメントには **設計上の誤記** が含まれている。  
> - 誤: `github.com/anthropic-ai/sdk-go`（Anthropic SDK）を使った Claude API 直接呼び出し  
> - 正: ユーザー自身の **OpenAI アカウント（Codex 経由）** を使う（ADR-004 で確定）  
>
> **[更新 2026-07-19]** 正式経路は **Codex CLI**（`internal/extraction/codex_cli_provider.go`）として
> **実装済み**。HTTP/OpenAI 直叩き provider は credential 非管理原則により deprecated
> （`ai-foundation-direction.md` / `codex-cli-integration.md` 参照）。
> section 4 のコード例（`callClaudeAPI` 等）は**歴史的な設計参考**であり、現実装とは異なる。
> pipeline 構造・validation・normalization の記述は引き続き有効。

> **[拡張 2026-07-06] 会話返信（reply）の追加**  
> チャットを「抽出専用」から「会話しながら抽出する」へ拡張した。AI 出力契約
> （`ExtractedFactResult`）に **`reply`（ユーザーへの会話返信）** を追加し、
> **`extracted_facts` は 0 件を許容**する（0 件の場合は `reply` 必須）。
> 非 fact の発言（質問・雑談）には `reply` で応答し、キャリアの話へ誘導する。
> fact 捏造禁止の原則は不変。§3 の型定義と §5 の prompt はこの契約が正。

> **[拡張 2026-07-19] 会話履歴（multi-turn）の導入**  
> conversation には直近の会話 transcript（`user:` / `assistant:` プレフィクス付き行、
> 末尾の `user:` 行が最新発言）を渡せる。プレフィクスが無い場合は単一発言として扱う。
> 会話中の既出 fact に詳細が追加された場合（clarification への回答など）、AI は
> **同じ `fact_id_hint` を再利用**して完全な更新版 fact を再出力する（重複でなく更新）。
> 同一 hint → 同一 `fact_id` → 既存の upsert が同一 id 置換するため、承認で fact が育つ。

AI が user との会話からキャリア fact を抽出し、patch proposal を生成するための Go ベースの仕様。

**Philosophy**: LLM を「structured JSON extraction provider」として扱い、extraction output を Go struct に deserialize → validate → normalize → patch generation する deterministic pipeline。AI provider は OpenAI Codex（ユーザー自身のアカウント）を使用する（ADR-004）。

---

## 1. Architecture Overview

```
User conversation
  ↓
Claude API (structured extraction prompt)
  ↓
Structured JSON response
  ↓
Go: JSON → ExtractedFact struct (deserialize)
  ↓
Go: Schema validation + normalization
  ↓
Go: Patch generation (fact_upsert operation)
  ↓
Patch YAML (user review 用)
  ↓
User review + approval
  ↓
Go: Patch apply (validate + write to facts/*.yaml)
  ↓
facts/*.yaml に append
```

**Key**: すべての validation・normalization・schema handling は Go side で実施。AI output は trusted input ではなく、untrusted JSON として扱う。

---

## 2. Go Package Structure

```
implementation/apps/api/
├── internal/
│   ├── extraction/
│   │   ├── extractor.go          # Main extraction orchestrator
│   │   ├── prompts.go            # Prompt templates
│   │   ├── models.go             # Extracted* struct definitions
│   │   ├── validator.go          # Schema validation
│   │   └── normalizer.go         # JSON → YAMLFact normalization
│   ├── patch/
│   │   ├── generator.go          # Fact → Patch conversion
│   │   └── applier.go            # Patch apply logic
│   └── workspace/
│       └── schema.go             # MVP YAML schema definitions
├── cmd/
│   └── server/
│       └── main.go               # API server (with /extract endpoint)
└── tests/
    ├── extraction_test.go
    ├── validator_test.go
    └── testdata/                 # fixtures
```

---

## 3. Core Go Types（Extraction Output）

### ExtractedFactResult (Claude API response)

```go
package extraction

// ExtractedFactResult は AI から返される構造化出力
type ExtractedFactResult struct {
	// Reply はユーザーへの会話返信。extracted_facts が 0 件の場合は必須
	//（非 fact の発言には reply で応答し、会話を継続する）。
	Reply              string          `json:"reply,omitempty"`
	ExtractedFacts     []ExtractedFact `json:"extracted_facts"` // 0 件を許容
	ExtractionQuality  ExtractionQuality `json:"extraction_quality"`
}

// ExtractedFact は 1 つ抽出された fact
type ExtractedFact struct {
	Type                   string                 `json:"type"` // experience | achievement | skill
	FactIDHint             string                 `json:"fact_id_hint"` // Semantic ID hint
	Summary                string                 `json:"summary"`
	Period                 *PeriodInfo            `json:"period,omitempty"` // For experience only
	Company                string                 `json:"company,omitempty"`
	Description            string                 `json:"description"`
	Confidence             string                 `json:"confidence"` // high | medium | low
	Details                map[string]interface{} `json:"details"` // Type-specific details
	ExtractionNotes        []string               `json:"extraction_notes"`
	ClarificationQuestions []string               `json:"clarification_questions"`
}

type PeriodInfo struct {
	Start string `json:"start"` // YYYY-MM or "unknown"
	End   string `json:"end"`
}

type ExtractionQuality struct {
	OverallConfidence        string `json:"overall_confidence"` // high | medium | low
	Completeness             string `json:"completeness"`
	NeedsClarificationCount  int    `json:"needs_clarification_count"`
	Summary                  string `json:"summary"`
}
```

### YAMLFact (MVP schema-compliant struct)

```go
package workspace

// YAMLFact は MVP YAML schema に準拠する fact
type YAMLFact struct {
	FactID      string            `yaml:"fact_id"`
	Type        string            `yaml:"type"` // experience | achievement | skill
	Status      string            `yaml:"status"` // confirmed | proposed
	Summary     string            `yaml:"summary"`
	Period      string            `yaml:"period,omitempty"` // For experience
	Company     string            `yaml:"company,omitempty"`
	Description string            `yaml:"description"`
	Confidence  string            `yaml:"confidence"`
	Source      string            `yaml:"source"`
	CreatedAt   string            `yaml:"created_at"`
	UpdatedAt   string            `yaml:"updated_at,omitempty"`
	Tags        []string          `yaml:"tags"`
	TechStack   []string          `yaml:"tech_stack,omitempty"`
	TeamSize    *int              `yaml:"team_size,omitempty"`
	Impact      map[string]interface{} `yaml:"impact,omitempty"`
	SourceDetail string           `yaml:"source_detail"`
}
```

---

## 4. Extraction Pipeline（Go Implementation）

### 4.1 Main Orchestrator

```go
package extraction

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropic-ai/sdk-go"
	"careerness/internal/workspace"
)

type Extractor struct {
	client *sdk.Client
}

func NewExtractor(apiKey string) *Extractor {
	return &Extractor{
		client: sdk.NewClient(apiKey),
	}
}

// ExtractFromConversation orchestrates the full extraction pipeline
func (e *Extractor) ExtractFromConversation(
	ctx context.Context,
	conversation string,
	sessionID string,
) (*ExtractionPipelineResult, error) {
	
	// Step 1: Call Claude API with structured prompt
	apiResult, err := e.callClaudeAPI(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("claude api call: %w", err)
	}

	// Step 2: Validate JSON structure
	if err := validateAPIResult(apiResult); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	// Step 3: Convert to YAMLFact structs (deserialize + normalize)
	yamlFacts := make([]*workspace.YAMLFact, 0)
	for _, extracted := range apiResult.ExtractedFacts {
		yamlFact, err := normalizeToYAMLFact(&extracted, sessionID)
		if err != nil {
			// Log error but continue processing other facts
			fmt.Printf("warning: failed to normalize fact %s: %v\n", extracted.FactIDHint, err)
			continue
		}
		yamlFacts = append(yamlFacts, yamlFact)
	}

	// Step 4: Validate against MVP schema
	for _, fact := range yamlFacts {
		if err := workspace.ValidateFact(fact); err != nil {
			return nil, fmt.Errorf("schema validation for %s: %w", fact.FactID, err)
		}
	}

	// Step 5: Generate patches
	patches, err := e.generatePatches(yamlFacts, sessionID)
	if err != nil {
		return nil, fmt.Errorf("patch generation: %w", err)
	}

	return &ExtractionPipelineResult{
		ExtractionQuality: apiResult.ExtractionQuality,
		YAMLFacts:         yamlFacts,
		Patches:           patches,
	}, nil
}

type ExtractionPipelineResult struct {
	ExtractionQuality ExtractionQuality
	YAMLFacts         []*workspace.YAMLFact
	Patches           []*workspace.Patch
}
```

### 4.2 Claude API Call

```go
// callClaudeAPI calls Claude with structured extraction prompt
func (e *Extractor) callClaudeAPI(
	ctx context.Context,
	conversation string,
) (*ExtractedFactResult, error) {
	
	// Get prompt template
	systemPrompt := getSystemPrompt()
	userPrompt := getUserPrompt(conversation)

	// Call API
	message, err := e.client.Messages.New(ctx, &sdk.MessageNewParams{
		Model: sdk.String("claude-opus-4-7"),
		Messages: []sdk.MessageParam{
			sdk.MessageParam{
				Role: "user",
				Content: userPrompt,
			},
		},
		System:    sdk.String(systemPrompt),
		MaxTokens: sdk.Int(2000),
	})
	if err != nil {
		return nil, err
	}

	// Extract JSON from response
	responseText := message.Content[0].Text

	// Parse JSON
	var result ExtractedFactResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	return &result, nil
}
```

### 4.3 Validator

```go
package extraction

import "fmt"

// validateAPIResult validates Claude API output structure
func validateAPIResult(result *ExtractedFactResult) error {
	if result == nil {
		return fmt.Errorf("nil result")
	}

	if len(result.ExtractedFacts) == 0 {
		return fmt.Errorf("no facts extracted")
	}

	for i, fact := range result.ExtractedFacts {
		// Validate required fields
		if fact.Type == "" {
			return fmt.Errorf("fact %d: missing type", i)
		}
		if fact.Type != "experience" && fact.Type != "achievement" && fact.Type != "skill" {
			return fmt.Errorf("fact %d: invalid type %q", i, fact.Type)
		}

		if fact.FactIDHint == "" {
			return fmt.Errorf("fact %d: missing fact_id_hint", i)
		}

		if fact.Summary == "" {
			return fmt.Errorf("fact %d: missing summary", i)
		}

		if fact.Confidence != "high" && fact.Confidence != "medium" && fact.Confidence != "low" {
			return fmt.Errorf("fact %d: invalid confidence %q", i, fact.Confidence)
		}
	}

	// Validate extraction quality
	if result.ExtractionQuality.OverallConfidence == "" {
		return fmt.Errorf("missing extraction_quality.overall_confidence")
	}

	return nil
}
```

### 4.4 Normalizer（JSON → YAML）

```go
package extraction

import (
	"fmt"
	"time"

	"careerness/internal/workspace"
)

// normalizeToYAMLFact converts ExtractedFact to YAMLFact
func normalizeToYAMLFact(
	extracted *ExtractedFact,
	sessionID string,
) (*workspace.YAMLFact, error) {

	now := time.Now().UTC().Format(time.RFC3339) + "Z"

	// Base fact (common to all types)
	fact := &workspace.YAMLFact{
		FactID:    generateFactID(extracted),
		Type:      extracted.Type,
		Status:    "proposed", // Always proposed from extraction
		Summary:   extracted.Summary,
		Confidence: extracted.Confidence,
		Source:    "conversation",
		CreatedAt: now,
		Tags:      inferTags(extracted),
		SourceDetail: formatSourceDetail(extracted),
	}

	// Type-specific normalization
	switch extracted.Type {
	case "experience":
		fact.Description = extracted.Description
		if extracted.Period != nil {
			fact.Period = formatPeriod(extracted.Period)
		}
		if extracted.Company != "" {
			fact.Company = extracted.Company
		}
		
		// Extract optional fields from Details
		if details, ok := extracted.Details.(map[string]interface{}); ok {
			if techStack, ok := details["tech_stack"].([]interface{}); ok {
				fact.TechStack = interfaceSliceToStrings(techStack)
			}
			if teamSize, ok := details["team_size"].(float64); ok {
				size := int(teamSize)
				fact.TeamSize = &size
			}
			if impact, ok := details["impact"].(map[string]interface{}); ok {
				fact.Impact = impact
			}
		}

	case "achievement":
		fact.Description = extracted.Description
		if details, ok := extracted.Details.(map[string]interface{}); ok {
			if impact, ok := details["impact"].(map[string]interface{}); ok {
				fact.Impact = impact
			}
		}

	case "skill":
		fact.Description = extracted.Description
		if details, ok := extracted.Details.(map[string]interface{}); ok {
			if proficiency, ok := details["proficiency"].(string); ok {
				// Could store in fact if schema extended
				_ = proficiency
			}
		}
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

// inferTags infers tags from extracted content
func inferTags(fact *ExtractedFact) []string {
	tags := make(map[string]bool)

	// Keyword-based inference
	fullText := fact.Summary + " " + fact.Description
	
	if containsAnyLower(fullText, []string{"backend", "api", "server", "database"}) {
		tags["backend"] = true
	}
	if containsAnyLower(fullText, []string{"frontend", "ui", "react", "javascript"}) {
		tags["frontend"] = true
	}
	if containsAnyLower(fullText, []string{"platform", "kubernetes", "devops", "infra"}) {
		tags["platform"] = true
	}
	if containsAnyLower(fullText, []string{"performance", "latency", "optimization"}) {
		tags["performance"] = true
	}
	if containsAnyLower(fullText, []string{"lead", "team", "coordination", "mentoring"}) {
		tags["leadership"] = true
	}

	// Convert to sorted slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
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

	return joinLines(lines)
}

// Helper functions
func formatPeriod(p *PeriodInfo) string {
	return fmt.Sprintf("%s to %s", p.Start, p.End)
}

func interfaceSliceToStrings(slice []interface{}) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func containsAnyLower(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(text), kw) {
			return true
		}
	}
	return false
}

func joinLines(lines []string) string {
	// ...
}
```

---

## 5. Prompt Templates（Go embedded resources）

### prompts.go

```go
package extraction

// getSystemPrompt returns the system prompt for conversational fact extraction
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

// getUserPrompt returns the user prompt with conversation
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
```

---

## 6. Patch Generation（Go）

```go
package patch

import (
	"fmt"
	"time"

	"careerness/internal/workspace"
)

// GeneratePatchFromFact creates a patch proposal from a YAMLFact
func GeneratePatchFromFact(
	fact *workspace.YAMLFact,
	sessionID string,
	index int,
) *workspace.Patch {

	now := time.Now().UTC().Format(time.RFC3339) + "Z"
	targetFile := fmt.Sprintf("facts/%ss.yaml", fact.Type) // experiences.yaml, etc.

	patch := &workspace.Patch{
		PatchID:     fmt.Sprintf("patch-%s", sessionID[6:]),
		Kind:        "fact_upsert",
		Status:      "proposed",
		CreatedAt:   now,
		CreatedBy:   "ai",
		WorkspaceID: "local-careervault",
		SessionID:   sessionID,
		Summary: fmt.Sprintf(
			"Extracted %s fact: %s\nConfidence: %s. Status: proposed.",
			fact.Type,
			fact.Summary,
			fact.Confidence,
		),
		Operations: []workspace.Operation{
			{
				OpID:   "op-001",
				Type:   "upsert_fact",
				Target: targetFile,
				FactID: fact.FactID,
				NewFact: fact,
			},
		},
	}

	return patch
}
```

---

## 7. API Endpoint

```go
// POST /api/v1/extract
// Request: {conversation: string}
// Response: {patches: [...], quality: {...}}

func (s *Server) extractFacts(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Conversation string `json:"conversation"`
		SessionID    string `json:"session_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Generate session ID if not provided
	if req.SessionID == "" {
		req.SessionID = fmt.Sprintf("sess_%s", shortID())
	}

	// Extract
	result, err := s.extractor.ExtractFromConversation(r.Context(), req.Conversation, req.SessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("extraction failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Response
	response := map[string]interface{}{
		"session_id":         req.SessionID,
		"extraction_quality": result.ExtractionQuality,
		"patches":            result.Patches,
		"preview": map[string]interface{}{
			"fact_count": len(result.YAMLFacts),
			"facts":      result.YAMLFacts,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

---

## 8. Validation（Go side）

### workspace/schema.go

```go
package workspace

import "fmt"

// ValidateFact validates a YAMLFact against MVP schema
func ValidateFact(fact *YAMLFact) error {
	// Required fields
	if fact.FactID == "" {
		return fmt.Errorf("missing fact_id")
	}
	if fact.Type == "" || (fact.Type != "experience" && fact.Type != "achievement" && fact.Type != "skill") {
		return fmt.Errorf("invalid type: %q", fact.Type)
	}
	if fact.Status == "" || (fact.Status != "confirmed" && fact.Status != "proposed") {
		return fmt.Errorf("invalid status: %q", fact.Status)
	}
	if fact.Summary == "" {
		return fmt.Errorf("missing summary")
	}
	if fact.Description == "" {
		return fmt.Errorf("missing description")
	}
	if fact.Confidence == "" || (fact.Confidence != "high" && fact.Confidence != "medium" && fact.Confidence != "low") {
		return fmt.Errorf("invalid confidence: %q", fact.Confidence)
	}
	if fact.Source == "" {
		return fmt.Errorf("missing source")
	}
	if fact.CreatedAt == "" {
		return fmt.Errorf("missing created_at")
	}

	// Extracted facts must be proposed
	if fact.Source == "conversation" && fact.Status != "proposed" {
		return fmt.Errorf("extracted fact must have status=proposed, got %q", fact.Status)
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

	// Tag validation (TODO: check against tags.yaml)
	// ...

	return nil
}
```

---

## 9. MVP Scope（Current）

### In Scope

- ✅ Claude API call with structured prompt
- ✅ JSON deserialize → Go struct
- ✅ Schema validation (required fields, type checks)
- ✅ Normalization (fact ID generation, tag inference)
- ✅ Patch generation (fact_upsert operation)
- ✅ API endpoint for extraction
- ✅ Error handling + logging

### Out of Scope（v2+）

- ❌ Multi-stage extraction (sequential refinement)
- ❌ Graph inference (relationship detection)
- ❌ Semantic deduplication
- ❌ Automated fact merging
- ❌ Complex confidence propagation
- ❌ Adaptive prompting based on facts

---

## 10. Testing Strategy

### Unit Tests

```go
// extraction_test.go
func TestNormalizeToYAMLFact(t *testing.T) {
	extracted := &ExtractedFact{
		Type: "experience",
		FactIDHint: "payment-platform",
		Summary: "Payment platform redesign",
		// ...
	}

	fact, err := normalizeToYAMLFact(extracted, "sess_test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fact.FactID != "fact-proj-payment-platform" {
		t.Errorf("expected fact-proj-payment-platform, got %q", fact.FactID)
	}

	if fact.Status != "proposed" {
		t.Errorf("extracted fact should be proposed, got %q", fact.Status)
	}
}

// validator_test.go
func TestValidateFact(t *testing.T) {
	tests := []struct{
		name string
		fact *YAMLFact
		valid bool
	}{
		{
			name: "valid experience",
			fact: &YAMLFact{
				FactID: "fact-proj-example",
				Type: "experience",
				Status: "proposed",
				Summary: "Example",
				Description: "...",
				Confidence: "high",
				Source: "conversation",
				CreatedAt: "2026-05-24T...",
				Company: "Example Corp",
				Period: "2022-04 to 2023-01",
			},
			valid: true,
		},
		{
			name: "missing company",
			fact: &YAMLFact{
				Type: "experience",
				// ... no company ...
			},
			valid: false,
		},
	}

	// ...
}
```

### Integration Tests

- Mock Claude API response
- Test full pipeline: JSON → validation → normalization → patch
- Verify YAML output matches schema

---

## 11. Error Handling

All errors must be:
1. **Logged**: `log.Printf("[extraction] error: %v", err)`
2. **Wrapped**: `fmt.Errorf("extraction: %w", err)`
3. **Handled gracefully**: User sees clear message
4. **Tracked**: Session ID preserved for debugging

Example:
```go
if err := validateAPIResult(result); err != nil {
	s.logger.Printf("[%s] validation failed: %v", sessionID, err)
	return &ErrorResponse{
		Code: "VALIDATION_FAILED",
		Message: "Extracted fact did not meet schema requirements",
		Details: err.Error(),
	}
}
```

---

## 12. Implementation Checklist

- [ ] Define `extraction.ExtractedFact`, `extraction.ExtractedFactResult` structs
- [ ] Define `workspace.YAMLFact` struct (MVP schema)
- [ ] Implement `callClaudeAPI` (Claude SDK call)
- [ ] Implement `validateAPIResult` (JSON validation)
- [ ] Implement `normalizeToYAMLFact` (JSON → YAML conversion)
- [ ] Implement `workspace.ValidateFact` (schema validation)
- [ ] Implement `patch.GeneratePatchFromFact`
- [ ] Implement `POST /api/v1/extract` endpoint
- [ ] Unit tests for normalizer, validator
- [ ] Integration test (API → patches)
- [ ] Example extraction runs with test data

---

## Next: File Structure Implementation

Follows this specification with Go packages and interfaces. Focus on:
1. **Clean interfaces**: extraction.Extractor, workspace.Schema
2. **Deterministic output**: Same input → same output
3. **Comprehensive logging**: Session ID everywhere
4. **Simple tests**: Table-driven, no mocks needed yet

**Goal**: "Extraction pipeline works" by end of week.
