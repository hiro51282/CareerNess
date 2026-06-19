package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"careerness/api/internal/session"
)

func TestPostAttach_RegistersSession(t *testing.T) {
	store := session.NewStore()
	root := t.TempDir()

	body := `{"workspace_id":"my-vault","workspace_root":"` + root + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspace/attach", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostAttach(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}
	sid := resp["session_id"]
	if sid == "" {
		t.Fatal("session_id が返らない")
	}

	att, ok := store.Get(sid)
	if !ok {
		t.Fatal("store に attachment が登録されていない")
	}
	if att.WorkspaceID != "my-vault" {
		t.Errorf("workspace_id = %q, want my-vault", att.WorkspaceID)
	}
	// root は正規化（EvalSymlinks）されるため文字列一致は保証しないが、空であってはならない。
	if att.WorkspaceRoot == "" {
		t.Error("workspace_root が空")
	}
}

func TestPostAttach_MissingRoot(t *testing.T) {
	store := session.NewStore()
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"workspace_id":"v"}`))
	rec := httptest.NewRecorder()

	PostAttach(store)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestPostAttach_NonexistentRoot(t *testing.T) {
	store := session.NewStore()
	body := `{"workspace_root":"/no/such/dir/really/unlikely"}`
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostAttach(store)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
