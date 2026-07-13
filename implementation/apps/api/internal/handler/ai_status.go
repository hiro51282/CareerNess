package handler

import (
	"net/http"
	"os"

	"careerness/api/internal/extraction"
)

// GetAIStatus は AI 実行環境の利用可否を返す（オンボーディング用・トークン消費なし）。
// UI は起動時と「再確認」時にこれを取得し、バッジ・案内バナーを表示する。
// アプリ内から codex login を起動する動線は Desktop Host 判断とセットで将来対応。
func GetAIStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET のみ受け付けます")
		return
	}

	type resp struct {
		Provider string `json:"provider"`
		Ready    bool   `json:"ready"`
		Detail   string `json:"detail"`
		Guidance string `json:"guidance,omitempty"`
	}

	switch os.Getenv("EXTRACTION_PROVIDER") {
	case "", "mock":
		writeJSON(w, http.StatusOK, resp{
			Provider: "mock",
			Ready:    true,
			Detail:   "Mock Provider（開発用）。実 AI は EXTRACTION_PROVIDER=codex-cli で有効化します。",
		})
	case "codex-cli":
		st := extraction.CodexCLIStatus(r.Context(), extraction.CodexCLIConfig{
			Bin:   os.Getenv("CODEX_CLI_BIN"),
			Model: os.Getenv("CODEX_CLI_MODEL"),
		})
		writeJSON(w, http.StatusOK, resp{Provider: "codex-cli", Ready: st.Ready, Detail: st.Detail, Guidance: st.Guidance})
	case "codex":
		ready := os.Getenv("OPENAI_API_KEY") != ""
		guidance := ""
		if !ready {
			guidance = "OPENAI_API_KEY が未設定です（この provider は deprecated。codex-cli の利用を推奨します）。"
		}
		writeJSON(w, http.StatusOK, resp{Provider: "codex", Ready: ready, Detail: "HTTP provider（deprecated）", Guidance: guidance})
	default:
		writeJSON(w, http.StatusOK, resp{
			Provider: os.Getenv("EXTRACTION_PROVIDER"),
			Ready:    false,
			Detail:   "未知の provider 設定",
			Guidance: "EXTRACTION_PROVIDER には mock または codex-cli を指定してください。",
		})
	}
}
