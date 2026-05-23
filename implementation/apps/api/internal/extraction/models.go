package extraction

// ExtractedFactResult は Claude API から返される構造化出力
type ExtractedFactResult struct {
	ExtractedFacts    []ExtractedFact   `json:"extracted_facts"`
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

// ExtractionPipelineResult は extraction pipeline の全結果
type ExtractionPipelineResult struct {
	ExtractionQuality ExtractionQuality `json:"extraction_quality"`
	YAMLFacts         []*YAMLFact       `json:"yaml_facts"`
	Patches           []*Patch          `json:"patches"`
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

// Patch は patch proposal
type Patch struct {
	PatchID     string      `yaml:"patch_id" json:"patch_id"`
	Kind        string      `yaml:"kind" json:"kind"` // fact_upsert
	Status      string      `yaml:"status" json:"status"` // proposed
	CreatedAt   string      `yaml:"created_at" json:"created_at"`
	CreatedBy   string      `yaml:"created_by" json:"created_by"` // ai
	WorkspaceID string      `yaml:"workspace_id" json:"workspace_id"`
	SessionID   string      `yaml:"session_id" json:"session_id"`
	Summary     string      `yaml:"summary" json:"summary"`
	Operations  []Operation `yaml:"operations" json:"operations"`
}

// Operation は patch 内の 1 つの操作
type Operation struct {
	OpID    string    `yaml:"op_id" json:"op_id"`
	Type    string    `yaml:"type" json:"type"` // upsert_fact
	Target  string    `yaml:"target" json:"target"` // facts/experiences.yaml
	FactID  string    `yaml:"fact_id" json:"fact_id"`
	NewFact *YAMLFact `yaml:"new_fact" json:"new_fact"`
}
