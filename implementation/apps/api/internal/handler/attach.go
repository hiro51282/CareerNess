package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"careerness/api/internal/session"
)

// PostAttach は workspace を session に attach する（ADR-006）。
//
// workspace_root を正規化（絶対化・symlink 解決）して実在ディレクトリであることを
// 確認し、session store に束縛を登録する。以降の apply はこの root 配下に
// 封じ込められ、apply 時に root をリクエストボディから受け取ることはない。
func PostAttach(store *session.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "POST のみ受け付けます")
			return
		}

		var req struct {
			SessionID     string `json:"session_id"`
			WorkspaceID   string `json:"workspace_id"`
			WorkspaceRoot string `json:"workspace_root"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "リクエストの解析に失敗しました")
			return
		}
		if strings.TrimSpace(req.WorkspaceRoot) == "" {
			writeError(w, http.StatusBadRequest, "workspace_root は必須です")
			return
		}

		realRoot, err := normalizeRoot(req.WorkspaceRoot)
		if err != nil {
			writeError(w, http.StatusBadRequest, "workspace_root が不正です: "+err.Error())
			return
		}

		sessionID := req.SessionID
		if sessionID == "" {
			sessionID = "sess-" + shortID()
		}
		workspaceID := req.WorkspaceID
		if workspaceID == "" {
			workspaceID = filepath.Base(realRoot)
		}

		store.Put(session.Attachment{
			SessionID:     sessionID,
			WorkspaceID:   workspaceID,
			WorkspaceRoot: realRoot,
		})

		writeJSON(w, http.StatusOK, map[string]string{
			"session_id":   sessionID,
			"workspace_id": workspaceID,
		})
	}
}

// normalizeRoot は workspace_root を絶対化・symlink 解決し、
// 実在するディレクトリであることを確認した実パスを返す。
//
// CodeQL Autofix（01c44a8）が追加した固定 allowlist `/workspaces` 配下チェックは
// ここで意図的に撤回している。CareerNESS は Local-first（ADR-001 / ADR-003）であり、
// workspace_root は「ユーザーが明示宣言した capability 境界」そのもので、任意の場所に
// 置かれ得る。固定 allowlist はこの設計を破壊する。書き込み時の root 外封じ込めは
// 後段の workspace.ResolveWithin（ADR-006）が root 基準で担保する。
// マルチユーザー/クラウド構成での扱いは Task3（認証）で再訪する。
func normalizeRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("ディレクトリが存在しません")
	}
	info, err := os.Stat(real)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("ディレクトリではありません")
	}
	return real, nil
}
