package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestDist は index.html と assets/app.js を持つ疑似 dist を作る。
func newTestDist(t *testing.T) string {
	t.Helper()
	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("<html>INDEX</html>"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dist, "assets"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dist, "assets", "app.js"), []byte("console.log('app')"), 0644); err != nil {
		t.Fatal(err)
	}
	return dist
}

func newSPARecorder(t *testing.T, dist, path string) *httptest.ResponseRecorder {
	t.Helper()
	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // API 到達の目印
	})
	h := NewSPAHandler(dist, api)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestSPAHandler_ServesIndexAtRoot(t *testing.T) {
	rec := newSPARecorder(t, newTestDist(t), "/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "INDEX") {
		t.Errorf("root は index.html を返すべき: code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestSPAHandler_ServesStaticAsset(t *testing.T) {
	rec := newSPARecorder(t, newTestDist(t), "/assets/app.js")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "console.log") {
		t.Errorf("実在アセットはそのまま配信すべき: code=%d", rec.Code)
	}
}

func TestSPAHandler_FallsBackToIndex(t *testing.T) {
	rec := newSPARecorder(t, newTestDist(t), "/some/client/route")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "INDEX") {
		t.Errorf("未知パスは index.html にフォールバックすべき: code=%d", rec.Code)
	}
}

func TestSPAHandler_RoutesAPI(t *testing.T) {
	if rec := newSPARecorder(t, newTestDist(t), "/api/v1/anything"); rec.Code != http.StatusTeapot {
		t.Errorf("/api/ は API へルーティングすべき: code=%d", rec.Code)
	}
	if rec := newSPARecorder(t, newTestDist(t), "/health"); rec.Code != http.StatusTeapot {
		t.Errorf("/health は API へルーティングすべき: code=%d", rec.Code)
	}
}

func TestSPAHandler_TraversalContained(t *testing.T) {
	// dist の外（親ディレクトリ）に秘密ファイルを置き、.. で届かないことを確認する。
	parent := t.TempDir()
	dist := filepath.Join(parent, "dist")
	if err := os.MkdirAll(dist, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("INDEX"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(parent, "secret.txt"), []byte("SECRET"), 0644); err != nil {
		t.Fatal(err)
	}

	rec := newSPARecorder(t, dist, "/../secret.txt")
	if strings.Contains(rec.Body.String(), "SECRET") {
		t.Fatal("パストラバーサルで dist 外が読めてしまった")
	}
}
