package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	patcher "careerness/api/internal/patch"
	"careerness/api/internal/session"
)

// PostApplyPatch は承認済み patch を attach 済み session の workspace へ適用する（ADR-006）。
//
// 書き込み先 root はリクエストボディからではなく session store の attachment から
// 導出する。session_id に対応する attachment が無ければ拒否する。Applier 入口で
// patch.Validate と workspace.ResolveWithin による封じ込めが行われる。
func PostApplyPatch(store *session.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
			return
		}

		var req struct {
			SessionID string         `json:"session_id"`
			Patch     *patcher.Patch `json:"patch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "リクエストの解析に失敗しました")
			return
		}
		if strings.TrimSpace(req.SessionID) == "" {
			writeError(w, http.StatusBadRequest, "session_id は必須です")
			return
		}
		if req.Patch == nil {
			writeError(w, http.StatusBadRequest, "patch は必須です")
			return
		}

		// 認可: session に対応する workspace attachment を引く。root はここから導出し、
		// リクエストボディの任意パスは信用しない。
		att, ok := store.Get(req.SessionID)
		if !ok {
			writeError(w, http.StatusForbidden, "この session には workspace が attach されていません")
			return
		}

		// workspace の取り違え防止: patch の workspace_id が attachment と一致すること。
		if req.Patch.WorkspaceID != att.WorkspaceID {
			writeError(w, http.StatusConflict, "patch の workspace_id が attach された workspace と一致しません")
			return
		}

		applier := patcher.NewApplier(att.WorkspaceRoot)
		result, err := applier.ApplyPatch(req.Patch)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, "適用に失敗しました: "+err.Error())
			return
		}

		if !result.Success() && result.AppliedCount == 0 {
			writeError(w, http.StatusInternalServerError, "すべての操作が失敗しました")
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}
