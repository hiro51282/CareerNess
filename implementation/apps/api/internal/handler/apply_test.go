package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"careerness/api/internal/extraction"
	"careerness/api/internal/patch"
	"careerness/api/internal/session"
)

// validFact は patch を組み立てるための最小 YAMLFact を返す。
func validFact(id string) *extraction.YAMLFact {
	return &extraction.YAMLFact{
		FactID:      id,
		Type:        "experience",
		Status:      "proposed",
		Summary:     "summary",
		Description: "description",
		Confidence:  "high",
		Source:      "conversation",
		CreatedAt:   "2026-01-01T00:00:00Z",
		Tags:        []string{"backend"},
		Company:     "Test Corp",
		Period:      "2022-01 to 2023-01",
	}
}

// attachedStore は tempDir を attach 済みの store と session_id / workspace_id を返す。
func attachedStore(t *testing.T) (*session.Store, string, string, string) {
	t.Helper()
	root := t.TempDir()
	store := session.NewStore()
	sid := "sess-apply"
	wid := "my-vault"
	store.Put(session.Attachment{SessionID: sid, WorkspaceID: wid, WorkspaceRoot: root})
	return store, sid, wid, root
}

// doApply は apply-patch ハンドラを呼び出す。
func doApply(store *session.Store, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/apply-patch", strings.NewReader(body))
	rec := httptest.NewRecorder()
	PostApplyPatch(store)(rec, req)
	return rec
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}

func TestApplyPatch_HappyPath(t *testing.T) {
	store, sid, wid, root := attachedStore(t)

	p := patch.BuildFactUpsert(validFact("fact-proj-test"), sid, 0)
	p.WorkspaceID = wid
	body := mustJSON(t, map[string]any{"session_id": sid, "patch": p})

	rec := doApply(store, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	// attachment root 配下に書き込まれていること
	if _, err := os.Stat(filepath.Join(root, "facts", "experiences.yaml")); err != nil {
		t.Fatalf("root 配下にファイルが作成されていない: %v", err)
	}
}

func TestApplyPatch_MissingSession(t *testing.T) {
	store, _, wid, _ := attachedStore(t)
	p := patch.BuildFactUpsert(validFact("fact-proj-test"), "sess-apply", 0)
	p.WorkspaceID = wid
	body := mustJSON(t, map[string]any{"patch": p}) // session_id 無し

	rec := doApply(store, body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestApplyPatch_UnattachedSession(t *testing.T) {
	store, _, wid, _ := attachedStore(t)
	p := patch.BuildFactUpsert(validFact("fact-proj-test"), "sess-unknown", 0)
	p.WorkspaceID = wid
	body := mustJSON(t, map[string]any{"session_id": "sess-unknown", "patch": p})

	rec := doApply(store, body)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestApplyPatch_WorkspaceMismatch(t *testing.T) {
	store, sid, _, _ := attachedStore(t)
	p := patch.BuildFactUpsert(validFact("fact-proj-test"), sid, 0)
	p.WorkspaceID = "other-vault" // attachment は my-vault
	body := mustJSON(t, map[string]any{"session_id": sid, "patch": p})

	rec := doApply(store, body)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
}

func TestApplyPatch_TraversalRejected(t *testing.T) {
	store, sid, wid, _ := attachedStore(t)
	p := patch.BuildFactUpsert(validFact("fact-proj-test"), sid, 0)
	p.WorkspaceID = wid
	p.Operations[0].Target = "../../etc/evil.yaml"
	body := mustJSON(t, map[string]any{"session_id": sid, "patch": p})

	rec := doApply(store, body)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

// TestApplyPatch_IgnoresBodyWorkspacePath は、ボディに workspace_path を混入させても
// 無視され、書き込みが attachment root に閉じることを検証する（回帰）。
func TestApplyPatch_IgnoresBodyWorkspacePath(t *testing.T) {
	store, sid, wid, root := attachedStore(t)
	bogus := t.TempDir() // 攻撃者が書かせたい別ディレクトリ

	p := patch.BuildFactUpsert(validFact("fact-proj-test"), sid, 0)
	p.WorkspaceID = wid
	body := mustJSON(t, map[string]any{"session_id": sid, "workspace_path": bogus, "patch": p})

	rec := doApply(store, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	// attachment root に書かれ、bogus には書かれていないこと
	if _, err := os.Stat(filepath.Join(root, "facts", "experiences.yaml")); err != nil {
		t.Fatalf("attachment root に書き込まれていない: %v", err)
	}
	if _, err := os.Stat(filepath.Join(bogus, "facts", "experiences.yaml")); err == nil {
		t.Fatal("ボディの workspace_path が使われてしまった")
	}
}
