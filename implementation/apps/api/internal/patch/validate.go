package patch

import (
	"errors"
	"fmt"
	"strings"
)

var allowedOpTypes = map[OpType]bool{
	OpCreateFile:              true,
	OpUpdateFile:              true,
	OpDeleteFile:              true,
	OpUpsertFact:              true,
	OpMarkFactStatus:          true,
	OpReplaceGeneratedProfile: true,
	OpReplaceGeneratedExport:  true,
	OpAppendHistoryRecord:     true,
}

var allowedConfidence = map[Confidence]bool{
	ConfidenceHigh:   true,
	ConfidenceMedium: true,
	ConfidenceLow:    true,
}

// Validate はパッチの構造的整合性とワークスペース境界ルールを検証する。
func Validate(p *Patch) error {
	if p.PatchID == "" {
		return errors.New("patch_id は必須です")
	}
	if p.WorkspaceID == "" {
		return errors.New("workspace_id は必須です")
	}
	if p.SessionID == "" {
		return errors.New("session_id は必須です")
	}
	if len(p.Operations) == 0 {
		return errors.New("operations が空です")
	}
	for i := range p.Operations {
		if err := validateOp(&p.Operations[i]); err != nil {
			return err
		}
	}
	return nil
}

func validateOp(op *Operation) error {
	if op.OpID == "" {
		return errors.New("op_id は必須です")
	}
	if !allowedOpTypes[op.Type] {
		return fmt.Errorf("op %s: 未知の operation type %q", op.OpID, op.Type)
	}
	if op.Target == "" {
		return fmt.Errorf("op %s: target は必須です", op.OpID)
	}
	if strings.Contains(op.Target, "..") {
		return fmt.Errorf("op %s: target にパストラバーサルが含まれています", op.OpID)
	}
	if strings.HasPrefix(op.Target, "/") {
		return fmt.Errorf("op %s: target は相対パスである必要があります", op.OpID)
	}
	if op.Rationale == "" {
		return fmt.Errorf("op %s: rationale は必須です", op.OpID)
	}
	if !allowedConfidence[op.Confidence] {
		return fmt.Errorf("op %s: confidence は high/medium/low のいずれかです", op.OpID)
	}
	if op.Type == OpUpsertFact || op.Type == OpMarkFactStatus {
		if op.EntityID == "" {
			return fmt.Errorf("op %s: fact 操作には entity_id が必要です", op.OpID)
		}
	}
	// delete は常に要確認
	if op.Type == OpDeleteFile {
		op.ReviewRequired = true
	}
	return nil
}
