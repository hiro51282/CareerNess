package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NewSPAHandler は SPA（ビルド済み dist）と API を単一 listener で配信する
// デスクトップ（Electron）用ハンドラを返す（ADR-008: BrowserWindow は Go 配信の
// localhost を読む構成。フロントと API が同一オリジンになり CORS 不要）。
//
// ルーティング:
//   - /api/ と /health → api ハンドラ
//   - 実在する静的ファイル → dist から配信
//   - それ以外 → index.html（SPA フォールバック）
func NewSPAHandler(distDir string, api http.Handler) http.Handler {
	fileServer := http.FileServer(http.Dir(distDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/health" {
			api.ServeHTTP(w, r)
			return
		}

		// "/" を先頭に付けて Clean することで ".." を封じ込めた rooted path にする。
		clean := filepath.Clean("/" + r.URL.Path)
		p := filepath.Join(distDir, clean)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		// 未知のパス（クライアントサイドルーティング等）は index.html へ。
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})
}
