package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListFiles(t *testing.T) {
	root := t.TempDir()
	mk := func(rel, content string) {
		t.Helper()
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	mk("facts/experiences.yaml", "- fact_id: a")
	mk("meta/workspace.yml", "name: vault")
	mk("README.md", "# vault")
	mk("notes.txt", "ignore me")            // 対象外拡張子
	mk("a/b/c/d/deep.yaml", "too: deep")    // 深さ 4 超（ファイル depth=5）

	files, err := ListFiles(root, 4)
	if err != nil {
		t.Fatalf("ListFiles error: %v", err)
	}
	if _, ok := files["facts/experiences.yaml"]; !ok {
		t.Error("facts/experiences.yaml が含まれない")
	}
	if _, ok := files["meta/workspace.yml"]; !ok {
		t.Error(".yml が含まれない")
	}
	if _, ok := files["README.md"]; !ok {
		t.Error(".md が含まれない")
	}
	if _, ok := files["notes.txt"]; ok {
		t.Error(".txt は対象外のはず")
	}
	if _, ok := files["a/b/c/d/deep.yaml"]; ok {
		t.Error("深さ超過のファイルは対象外のはず")
	}
	if files["facts/experiences.yaml"] != "- fact_id: a" {
		t.Errorf("本文が一致しない: %q", files["facts/experiences.yaml"])
	}
}

func TestListFiles_SkipsSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.yaml"), []byte("SECRET"), 0644); err != nil {
		t.Fatal(err)
	}
	// root 内の symlink（ファイル・ディレクトリ両方）
	if err := os.Symlink(filepath.Join(outside, "secret.yaml"), filepath.Join(root, "link.yaml")); err != nil {
		t.Skipf("symlink 不可の環境: %v", err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linkdir")); err != nil {
		t.Skipf("symlink 不可の環境: %v", err)
	}

	files, err := ListFiles(root, 4)
	if err != nil {
		t.Fatalf("ListFiles error: %v", err)
	}
	for p, c := range files {
		if c == "SECRET" {
			t.Errorf("symlink 経由で root 外が読めてしまった: %s", p)
		}
	}
}
