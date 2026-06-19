package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveWithin_Normal は root 配下の正当な target が root 内に解決されることを検証する。
func TestResolveWithin_Normal(t *testing.T) {
	root := t.TempDir()

	cases := []string{
		"facts/experiences.yaml",
		"a/b/c.yaml",
		"meta/workspace.yaml",
	}
	for _, target := range cases {
		got, err := ResolveWithin(root, target)
		if err != nil {
			t.Fatalf("ResolveWithin(%q) で予期しないエラー: %v", target, err)
		}
		// 解決結果が root 配下であること
		rel, err := filepath.Rel(root, got)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			t.Errorf("target %q が root 外に解決された: %q", target, got)
		}
	}
}

// TestResolveWithin_Rejects は境界突破・不正入力を拒否することを検証する。
func TestResolveWithin_Rejects(t *testing.T) {
	root := t.TempDir()

	cases := map[string]string{
		"空 target":             "",
		"絶対パス":                 "/etc/passwd",
		"先頭トラバーサル":            "../../etc/passwd",
		"途中トラバーサル":            "facts/../../x.yaml",
		"root 内に留まる .. も拒否": "facts/../experiences.yaml",
	}
	for name, target := range cases {
		if _, err := ResolveWithin(root, target); err == nil {
			t.Errorf("%s: target %q は拒否されるべき", name, target)
		}
	}
}

// TestResolveWithin_EmptyRoot は root 未指定を拒否する。
func TestResolveWithin_EmptyRoot(t *testing.T) {
	if _, err := ResolveWithin("", "facts/x.yaml"); err == nil {
		t.Error("空 root は拒否されるべき")
	}
}

// TestResolveWithin_NonexistentRoot は存在しない root を拒否する。
func TestResolveWithin_NonexistentRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "does-not-exist")
	if _, err := ResolveWithin(root, "facts/x.yaml"); err == nil {
		t.Error("存在しない root は拒否されるべき")
	}
}

// TestResolveWithin_SymlinkEscape は root 内の symlink が外部を指す場合に遮断することを検証する。
func TestResolveWithin_SymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir() // root の外

	// root/out -> outside という symlink を張る
	link := filepath.Join(root, "out")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink を作成できない環境: %v", err)
	}

	if _, err := ResolveWithin(root, "out/secret.yaml"); err == nil {
		t.Error("symlink 経由の root 外書き込みは遮断されるべき")
	}
}
