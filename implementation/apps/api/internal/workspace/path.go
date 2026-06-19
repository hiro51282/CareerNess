package workspace

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// ResolveWithin は patch operation の相対 target を、attach された root 配下の
// 絶対パスへ封じ込めて解決する（ADR-006 / ADR-001）。
//
// 文字列パターン検査ではなく filepath.Rel ベースの封じ込めで判定し、
// さらに root 内に存在する祖先の symlink を評価して、root 外へ抜ける
// シンボリックリンクも遮断する。root 外を指す target はエラーを返す。
func ResolveWithin(root, target string) (string, error) {
	if strings.TrimSpace(root) == "" {
		return "", fmt.Errorf("workspace root が空です")
	}
	if strings.TrimSpace(target) == "" {
		return "", fmt.Errorf("target が空です")
	}
	if filepath.IsAbs(target) {
		return "", fmt.Errorf("target は相対パスである必要があります: %q", target)
	}
	// Go 標準の局所性判定で、root を脱出し得る target を早期に弾く。
	// 絶対パス・空・".." での脱出・（Windows の）予約名を一括で拒否し、
	// CodeQL go/path-injection も sanitizer として認識できる形にする。
	if !filepath.IsLocal(target) {
		return "", fmt.Errorf("target は workspace 内の相対パスである必要があります: %q", target)
	}
	// IsLocal は root 内に留まる ".." を許容するが、本実装では曖昧さを避けるため
	// 明示的な ".." セグメントはより厳格に拒否する（挙動は従来どおり）。
	if slices.Contains(strings.Split(filepath.ToSlash(target), "/"), "..") {
		return "", fmt.Errorf("target にパストラバーサルが含まれています: %q", target)
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("root の絶対パス化に失敗: %w", err)
	}
	// root の実体（symlink 解決済み）を封じ込めの基準にする。
	// 存在しない root は attach 不能とみなしてエラーにする。
	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return "", fmt.Errorf("workspace root を解決できません: %w", err)
	}

	joined := filepath.Join(realRoot, target)

	// 封じ込めを、sink に渡る変数（joined）そのもので字句的に検査する。
	// ".." セグメントを既に拒否しているため joined は realRoot 配下に収まるが、
	// この HasPrefix ガードでその不変条件を joined 自身に対して明示する。
	// これにより静的解析（CodeQL go/path-injection）も joined を sanitize 済みと
	// 認識できる（変数を跨いだガードでは barrier と見なされないため）。
	prefix := realRoot + string(filepath.Separator)
	if joined != realRoot && !strings.HasPrefix(joined, prefix) {
		return "", fmt.Errorf("target が workspace 外を指しています: %q", target)
	}

	// symlink 経由の脱出を遮断する: 書き込み先は未存在のことがあるため、
	// 存在する祖先まで遡って実体解決し、その実体で再封じ込めを判定する。
	checkPath := resolveExistingAncestor(joined)
	rel, err := filepath.Rel(realRoot, checkPath)
	if err != nil {
		return "", fmt.Errorf("封じ込め判定に失敗: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("target が workspace 外を指しています (symlink): %q", target)
	}

	return joined, nil
}

// resolveExistingAncestor は p の存在する最も深い祖先を symlink 解決し、
// 未存在の残り部分を足し戻して返す。どの祖先も解決できなければ p を返す。
func resolveExistingAncestor(p string) string {
	cur := p
	for {
		if resolved, err := filepath.EvalSymlinks(cur); err == nil {
			suffix := strings.TrimPrefix(p, cur)
			return filepath.Join(resolved, suffix)
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return p
		}
		cur = parent
	}
}
