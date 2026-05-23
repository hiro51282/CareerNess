package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"careerness/api/internal/extraction"
)

// PostExtract orchestrates fact extraction from conversation.
// Returns patch proposals ready for user review.
func PostExtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
		return
	}

	var req struct {
		Conversation string `json:"conversation"`
		SessionID    string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの解析に失敗しました")
		return
	}

	if strings.TrimSpace(req.Conversation) == "" {
		writeError(w, http.StatusBadRequest, "conversation は必須です")
		return
	}

	if req.SessionID == "" {
		req.SessionID = "sess-" + shortID()
	}

	// Extract facts from conversation
	service := extraction.NewExtractionService(extraction.NewMockExtractionProvider())
	result, err := service.ExtractFromConversation(r.Context(), req.Conversation, req.SessionID)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "抽出に失敗しました: "+err.Error())
		return
	}

	if len(result.Patches) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "conversation から事実を抽出できませんでした")
		return
	}

	response := map[string]interface{}{
		"patches":              result.Patches,
		"extraction_quality":   result.ExtractionQuality,
		"yaml_facts":           result.YAMLFacts,
	}
	writeJSON(w, http.StatusOK, response)
}
