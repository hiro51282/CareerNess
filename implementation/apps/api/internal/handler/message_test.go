package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPostMessage_DefaultMock は、既定（Mock）で /message が ExtractionService 経由の
// 単一 patch proposal を 200 で返すことを検証する（PR-A の一本化・応答形不変）。
func TestPostMessage_DefaultMock(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "") // 既定 = mock

	body := `{"session_id":"sess-x","workspace_id":"my-vault","message":"2023年にABCで決済基盤をGoへ移行した"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/message", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Reply   string `json:"reply"`
		Patches []struct {
			WorkspaceID string `json:"workspace_id"`
			Operations  []struct {
				Type   string `json:"type"`
				Target string `json:"target"`
			} `json:"operations"`
		} `json:"patches"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}
	if resp.Reply == "" {
		t.Error("reply が空")
	}
	// 既定 Mock は単一 fact のため patches は 1 件（契約は複数対応）。
	if len(resp.Patches) != 1 {
		t.Fatalf("patches = %d, want 1", len(resp.Patches))
	}
	p0 := resp.Patches[0]
	if p0.WorkspaceID != "my-vault" {
		t.Errorf("workspace_id = %q, want my-vault", p0.WorkspaceID)
	}
	if len(p0.Operations) != 1 || p0.Operations[0].Type != "upsert_fact" {
		t.Errorf("operations が想定外: %+v", p0.Operations)
	}
	if p0.Operations[0].Target != "facts/experiences.yaml" {
		t.Errorf("target = %q, want facts/experiences.yaml", p0.Operations[0].Target)
	}
}

// TestPostMessage_ConversationalTurn は非 fact の発言（疑問形）が 200 で
// reply のみ（patches 0 件）を返すことを検証する（B: 自由対話）。
func TestPostMessage_ConversationalTurn(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "") // 既定 = mock

	body := `{"session_id":"sess-x","message":"他にどんな情報を書けばいいですか？"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/message", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostMessage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Reply   string           `json:"reply"`
		Patches []map[string]any `json:"patches"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}
	if resp.Reply == "" {
		t.Error("会話返信 reply が空")
	}
	if len(resp.Patches) != 0 {
		t.Errorf("patches = %d, want 0（非 fact の発言）", len(resp.Patches))
	}
}

func TestPostMessage_EmptyMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/message", strings.NewReader(`{"message":"  "}`))
	rec := httptest.NewRecorder()

	PostMessage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
