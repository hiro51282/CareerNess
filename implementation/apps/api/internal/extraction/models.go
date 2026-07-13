package extraction

// ExtractedFactResult は AI から返される構造化出力
type ExtractedFactResult struct {
	// Reply はユーザーへの会話返信。extracted_facts が 0 件の場合は必須
	//（非 fact の発言には reply で応答し、会話を継続する。extraction-specification.md）。
	Reply             string            `json:"reply,omitempty"`
	ExtractedFacts    []ExtractedFact   `json:"extracted_facts"` // 0 件を許容
	ExtractionQuality ExtractionQuality `json:"extraction_quality"`
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

// PeriodInfo は期間情報
type PeriodInfo struct {
	Start string `json:"start"` // YYYY-MM or "unknown"
	End   string `json:"end"`
}

// ExtractionQuality は抽出品質の指標
type ExtractionQuality struct {
	OverallConfidence       string `json:"overall_confidence"` // high | medium | low
	Completeness            string `json:"completeness"`       // high | medium | low
	NeedsClarificationCount int    `json:"needs_clarification_count"`
	Summary                 string `json:"summary"`
}

// ExtractionPipelineResult は extraction pipeline の全結果。
// patch 生成は patch パッケージ（patch.BuildFactUpsert）の責務であり、
// import 循環を避けるため本パッケージは patch を持たない。
type ExtractionPipelineResult struct {
	// Reply は AI の会話返信（provider の reply を透過）。facts 0 件時の応答に使う。
	Reply             string            `json:"reply,omitempty"`
	ExtractionQuality ExtractionQuality `json:"extraction_quality"`
	YAMLFacts         []*YAMLFact       `json:"yaml_facts"`
}

// YAMLFact は MVP YAML schema に準拠する fact
type YAMLFact struct {
	FactID       string                 `yaml:"fact_id" json:"fact_id"`
	Type         string                 `yaml:"type" json:"type"` // experience | achievement | skill
	Status       string                 `yaml:"status" json:"status"` // confirmed | proposed
	Summary      string                 `yaml:"summary" json:"summary"`
	Period       string                 `yaml:"period,omitempty" json:"period,omitempty"` // For experience
	Company      string                 `yaml:"company,omitempty" json:"company,omitempty"`
	Description  string                 `yaml:"description" json:"description"`
	// Phase 1 Minimal Fix（docs/implementation/workspace/fact-schema.md, 2026-05-31 確定）:
	// action/decision を空で先行追加し、Phase 2 拡張時の confirmed facts マイグレーションを回避する。
	// 空の場合は表示側で description にフォールバックする（impact.summary は既存の Impact map で表現可能）。
	Action       string                 `yaml:"action,omitempty" json:"action,omitempty"`     // 何をしたか（空なら description 表示）
	Decision     string                 `yaml:"decision,omitempty" json:"decision,omitempty"` // 何を判断したか
	Confidence   string                 `yaml:"confidence" json:"confidence"`
	Source       string                 `yaml:"source" json:"source"`
	CreatedAt    string                 `yaml:"created_at" json:"created_at"`
	UpdatedAt    string                 `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
	Tags         []string               `yaml:"tags" json:"tags"`
	TechStack    []string               `yaml:"tech_stack,omitempty" json:"tech_stack,omitempty"`
	TeamSize     *int                   `yaml:"team_size,omitempty" json:"team_size,omitempty"`
	Impact       map[string]interface{} `yaml:"impact,omitempty" json:"impact,omitempty"`
	SourceDetail string                 `yaml:"source_detail" json:"source_detail"`
}
