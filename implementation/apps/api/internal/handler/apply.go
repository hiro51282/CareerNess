package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	patcher "careerness/api/internal/patch"
)

// PostApplyPatch applies an approved patch to the workspace.
func PostApplyPatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
		return
	}

	var req struct {
		Patch         *patcher.Patch `json:"patch"`
		WorkspacePath string         `json:"workspace_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの解析に失敗しました")
		return
	}

	if req.Patch == nil {
		writeError(w, http.StatusBadRequest, "patch は必須です")
		return
	}

	if strings.TrimSpace(req.WorkspacePath) == "" {
		writeError(w, http.StatusBadRequest, "workspace_path は必須です")
		return
	}

	// Path safety: reject .. to prevent traversal
	cleanPath := filepath.Clean(req.WorkspacePath)
	if strings.Contains(cleanPath, "..") {
		writeError(w, http.StatusBadRequest, "invalid workspace_path: path traversal not allowed")
		return
	}

	if len(req.Patch.Operations) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "patch に操作がありません")
		return
	}

	// Apply patch
	applier := patcher.NewApplier(cleanPath)
	result, err := applier.ApplyPatch(req.Patch)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "適用に失敗しました: "+err.Error())
		return
	}

	if !result.Success() && result.AppliedCount == 0 {
		writeError(w, http.StatusInternalServerError, "すべての操作が失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, result)
}
