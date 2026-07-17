package handler

import (
	"net/http"
	"strings"

	"careerness/api/internal/session"
	"careerness/api/internal/workspace"
)

// GetWorkspaceFiles は attach 済み session の workspace から YAML/MD ファイルを
// 読み取り専用で列挙して返す（desktop モードの表示用。Workspace Gateway の
// read-only listing）。root は session store から導出し、ボディ/クエリの
// パスは一切受け取らない（ADR-006 と同じ原則）。
func GetWorkspaceFiles(store *session.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "GET のみ受け付けます")
			return
		}
		sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
		if sessionID == "" {
			writeError(w, http.StatusBadRequest, "session_id は必須です")
			return
		}
		att, ok := store.Get(sessionID)
		if !ok {
			writeError(w, http.StatusForbidden, "この session には workspace が attach されていません")
			return
		}

		files, err := workspace.ListFiles(att.WorkspaceRoot, 4)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "workspace の読み取りに失敗しました: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"workspace_id": att.WorkspaceID,
			"files":        files,
		})
	}
}
