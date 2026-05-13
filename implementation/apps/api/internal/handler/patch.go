package handler

import (
	"encoding/json"
	"net/http"

	"careerness/api/internal/patch"
)

// PostValidatePatch はパッチ構造を検証して結果を返す。
func PostValidatePatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
		return
	}

	var p patch.Patch
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "パッチの解析に失敗しました")
		return
	}

	if err := patch.Validate(&p); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "valid"})
}
