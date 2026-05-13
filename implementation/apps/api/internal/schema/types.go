package schema

// FactStatus は career fact の確認状態を表す。
type FactStatus string

const (
	FactStatusProposed  FactStatus = "proposed"
	FactStatusInferred  FactStatus = "inferred"
	FactStatusConfirmed FactStatus = "confirmed"
	FactStatusRejected  FactStatus = "rejected"
)

// Fact はユーザーのキャリアに関する構造化された事実を保持する。
type Fact struct {
	ProjectID string     `yaml:"project_id" json:"project_id"`
	Period    string     `yaml:"period,omitempty" json:"period,omitempty"`
	RoleHints []string   `yaml:"role_hints,omitempty" json:"role_hints,omitempty"`
	Facts     []string   `yaml:"facts" json:"facts"`
	Status    FactStatus `yaml:"status" json:"status"`
}

// Profile は特定の audience / role 向けに facts を再構成した見せ方。
type Profile struct {
	ID          string   `yaml:"id" json:"id"`
	TargetRole  string   `yaml:"target_role" json:"target_role"`
	SourceFacts []string `yaml:"source_facts" json:"source_facts"`
	Summary     string   `yaml:"summary" json:"summary"`
	Emphasis    []string `yaml:"emphasis" json:"emphasis"`
}

// Export は提出・共有のために整形された最終出力。
type Export struct {
	Format   string `yaml:"format" json:"format"`
	Profile  string `yaml:"profile" json:"profile"`
	Audience string `yaml:"audience" json:"audience"`
	Output   string `yaml:"output" json:"output"`
}
