export type PatchStatus = 'proposed' | 'approved' | 'applied' | 'rejected'
export type OpType =
  | 'create_file'
  | 'update_file'
  | 'delete_file'
  | 'upsert_fact'
  | 'mark_fact_status'
  | 'replace_generated_profile'
  | 'replace_generated_export'
  | 'append_history_record'
export type Confidence = 'high' | 'medium' | 'low'
export type FactStatus = 'proposed' | 'inferred' | 'confirmed' | 'rejected'

export interface ChangeRecord {
  before: unknown
  after: unknown
}

export interface Operation {
  op_id: string
  type: OpType
  target: string
  entity_id?: string
  change: ChangeRecord
  rationale: string
  confidence: Confidence
  review_required?: boolean
  fact_status_after?: FactStatus
}

export interface Patch {
  patch_id: string
  workspace_id: string
  session_id: string
  created_at: string
  created_by: string
  kind: string
  summary: string
  status: PatchStatus
  operations: Operation[]
}

export interface MessageResponse {
  reply: string
  patches: Patch[]
}

// ChatTurn は会話履歴の 1 ターン（message API へ同送する直近文脈）。
export interface ChatTurn {
  role: 'user' | 'ai'
  text: string
}
