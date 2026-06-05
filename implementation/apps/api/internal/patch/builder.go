package patch

import (
	"fmt"
	"time"

	"careerness/api/internal/extraction"
)

// BuildFactUpsert は抽出された YAMLFact から upsert_fact パッチ提案を生成する。
//
// docs/implementation/ai/ai-patch-format.md の Patch Envelope に準拠する：
//   - 1 patch = 1 semantic change（fact 追加）
//   - operation は change{before, after} を持ち、after に fact 本体を格納する
//   - fact 操作なので entity_id / fact_status_after / review_required を設定する
//   - rationale / confidence は必須に近いため常に埋める
func BuildFactUpsert(fact *extraction.YAMLFact, sessionID string, idx int) *Patch {
	now := time.Now().UTC().Format(time.RFC3339)
	targetFile := fmt.Sprintf("facts/%ss.yaml", fact.Type)
	patchID := fmt.Sprintf("patch-%s-%03d", now[0:10], idx+1)

	rationale := fmt.Sprintf(
		"会話から抽出した %s fact「%s」です。内容を確認してください。",
		fact.Type, fact.Summary,
	)

	return &Patch{
		PatchID:     patchID,
		WorkspaceID: "local-careervault",
		SessionID:   sessionID,
		CreatedAt:   now,
		CreatedBy:   "ai",
		Kind:        "workspace_patch",
		Summary:     fmt.Sprintf("会話から %s fact を抽出しました: %s", fact.Type, fact.Summary),
		Status:      StatusProposed,
		Operations: []Operation{
			{
				OpID:     "op-001",
				Type:     OpUpsertFact,
				Target:   targetFile,
				EntityID: fact.FactID,
				Change: ChangeRecord{
					Before: nil,
					After:  fact,
				},
				Rationale:       rationale,
				Confidence:      mapConfidence(fact.Confidence),
				FactStatusAfter: fact.Status,
				ReviewRequired:  true,
			},
		},
	}
}

// mapConfidence は fact の文字列 confidence を patch の Confidence 型に対応づける。
// 未知の値は medium に丸める（validation を必ず通すため）。
func mapConfidence(c string) Confidence {
	switch Confidence(c) {
	case ConfidenceHigh:
		return ConfidenceHigh
	case ConfidenceLow:
		return ConfidenceLow
	default:
		return ConfidenceMedium
	}
}
