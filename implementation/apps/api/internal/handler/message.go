package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"careerness/api/internal/ai"
)

// PostMessage は会話メッセージを受け取り、patch proposal を返す。
func PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
		return
	}

	var req ai.ProposeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの解析に失敗しました")
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		writeError(w, http.StatusBadRequest, "message は必須です")
		return
	}
	if req.WorkspaceID == "" {
		req.WorkspaceID = "local-careervault"
	}
	if req.SessionID == "" {
		req.SessionID = "sess-" + shortID()
	}

	result := ai.Propose(&req)
	writeJSON(w, http.StatusOK, result)
}
