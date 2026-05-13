package patch

// Status はパッチのライフサイクル状態。
type Status string

const (
	StatusProposed Status = "proposed"
	StatusApproved Status = "approved"
	StatusApplied  Status = "applied"
	StatusRejected Status = "rejected"
)

// OpType はワークスペース変更の種別。
type OpType string

const (
	OpCreateFile              OpType = "create_file"
	OpUpdateFile              OpType = "update_file"
	OpDeleteFile              OpType = "delete_file"
	OpUpsertFact              OpType = "upsert_fact"
	OpMarkFactStatus          OpType = "mark_fact_status"
	OpReplaceGeneratedProfile OpType = "replace_generated_profile"
	OpReplaceGeneratedExport  OpType = "replace_generated_export"
	OpAppendHistoryRecord     OpType = "append_history_record"
)

// Confidence は AI が生成した変更の確信度。
type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
)

// ChangeRecord は 1 operation の変更前後の状態。
type ChangeRecord struct {
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

// Operation はパッチ内の 1 つのアトミック変更。
type Operation struct {
	OpID            string       `json:"op_id"`
	Type            OpType       `json:"type"`
	Target          string       `json:"target"`
	EntityID        string       `json:"entity_id,omitempty"`
	Change          ChangeRecord `json:"change"`
	Rationale       string       `json:"rationale"`
	Confidence      Confidence   `json:"confidence"`
	ReviewRequired  bool         `json:"review_required,omitempty"`
	FactStatusAfter string       `json:"fact_status_after,omitempty"`
}

// Patch は 1 会話ターンの AI 提案による workspace 変更セット。
type Patch struct {
	PatchID     string      `json:"patch_id"`
	WorkspaceID string      `json:"workspace_id"`
	SessionID   string      `json:"session_id"`
	CreatedAt   string      `json:"created_at"`
	CreatedBy   string      `json:"created_by"`
	Kind        string      `json:"kind"`
	Summary     string      `json:"summary"`
	Status      Status      `json:"status"`
	Operations  []Operation `json:"operations"`
}
