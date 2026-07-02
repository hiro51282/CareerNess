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
		Reply string `json:"reply"`
		Patch *struct {
			WorkspaceID string `json:"workspace_id"`
			Operations  []struct {
				Type   string `json:"type"`
				Target string `json:"target"`
			} `json:"operations"`
		} `json:"patch"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}
	if resp.Reply == "" {
		t.Error("reply が空")
	}
	if resp.Patch == nil {
		t.Fatal("patch が nil（応答形は単一 patch のまま維持される）")
	}
	if resp.Patch.WorkspaceID != "my-vault" {
		t.Errorf("workspace_id = %q, want my-vault", resp.Patch.WorkspaceID)
	}
	if len(resp.Patch.Operations) != 1 || resp.Patch.Operations[0].Type != "upsert_fact" {
		t.Errorf("operations が想定外: %+v", resp.Patch.Operations)
	}
	if resp.Patch.Operations[0].Target != "facts/experiences.yaml" {
		t.Errorf("target = %q, want facts/experiences.yaml", resp.Patch.Operations[0].Target)
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
