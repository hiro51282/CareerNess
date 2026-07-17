package workspace

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ListFiles は root 配下の YAML / Markdown ファイルを読み取り専用で列挙し、
// 相対パス → 本文 の map を返す（Workspace Gateway の read-only listing。
// desktop モードでブラウザ FS 読み取りの代替として使う）。
//
// ブラウザ側 readWorkspaceFiles と同じ規約に揃える:
//   - 対象拡張子は .yaml / .yml / .md
//   - 深さは maxDepth まで（root 直下 = 1）
//   - symlink は辿らない（root 外への脱出を防ぐ）
//   - 1 ファイル 1MB 超はスキップ（表示用途のため）
func ListFiles(root string, maxDepth int) (map[string]string, error) {
	const maxFileSize = 1 << 20 // 1MB

	files := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}

		depth := len(strings.Split(filepath.ToSlash(rel), "/"))

		// symlink はファイル/ディレクトリとも扱わない（root 外脱出の防止）。
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			if depth >= maxDepth {
				return fs.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".yaml" && ext != ".yml" && ext != ".md" {
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil || info.Size() > maxFileSize {
			return nil
		}
		content, readErr := os.ReadFile(p)
		if readErr != nil {
			return nil // 読めないファイルはスキップ（表示用途）
		}
		files[filepath.ToSlash(rel)] = string(content)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
