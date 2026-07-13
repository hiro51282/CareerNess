package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getStatus(t *testing.T) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/status", nil)
	rec := httptest.NewRecorder()
	GetAIStatus(rec, req)
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}
	return rec.Code, body
}

func TestGetAIStatus_Mock(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "")
	code, body := getStatus(t)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if body["provider"] != "mock" || body["ready"] != true {
		t.Errorf("mock は ready=true であるべき: %v", body)
	}
}

func TestGetAIStatus_CodexCLIBinMissing(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "codex-cli")
	t.Setenv("CODEX_CLI_BIN", "/nonexistent/codex-xyz")
	code, body := getStatus(t)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if body["provider"] != "codex-cli" || body["ready"] != false {
		t.Errorf("bin 不在は ready=false であるべき: %v", body)
	}
	if g, _ := body["guidance"].(string); g == "" {
		t.Error("guidance が空")
	}
}

func TestGetAIStatus_UnknownProvider(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "bogus")
	code, body := getStatus(t)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if body["ready"] != false {
		t.Errorf("未知 provider は ready=false: %v", body)
	}
}
