package extraction

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeLoginCLI は `codex login status` を模した bash スクリプトを生成する。
func fakeLoginCLI(t *testing.T, exitCode int, output string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "fakecodex.sh")
	script := fmt.Sprintf("#!/usr/bin/env bash\necho %q\nexit %d\n", output, exitCode)
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCodexCLIStatus_BinMissing(t *testing.T) {
	st := CodexCLIStatus(context.Background(), CodexCLIConfig{Bin: "/nonexistent/codex-xyz"})
	if st.Ready {
		t.Error("バイナリ不在は ready=false であるべき")
	}
	if !strings.Contains(st.Guidance, "インストール") {
		t.Errorf("guidance にインストール案内が無い: %q", st.Guidance)
	}
}

func TestCodexCLIStatus_LoggedIn(t *testing.T) {
	bin := fakeLoginCLI(t, 0, "Logged in using ChatGPT")
	st := CodexCLIStatus(context.Background(), CodexCLIConfig{Bin: bin, Model: "gpt-5.4"})
	if !st.Ready {
		t.Errorf("ログイン済みは ready=true であるべき: %+v", st)
	}
	if !strings.Contains(st.Detail, "Logged in") {
		t.Errorf("detail = %q", st.Detail)
	}
	if st.Guidance != "" {
		t.Errorf("model 設定済みなら guidance は空であるべき: %q", st.Guidance)
	}
}

func TestCodexCLIStatus_NotLoggedIn(t *testing.T) {
	bin := fakeLoginCLI(t, 1, "Not logged in")
	st := CodexCLIStatus(context.Background(), CodexCLIConfig{Bin: bin, Model: "gpt-5.4"})
	if st.Ready {
		t.Error("未ログインは ready=false であるべき")
	}
	if !strings.Contains(st.Guidance, "codex login") {
		t.Errorf("guidance に codex login の案内が無い: %q", st.Guidance)
	}
}

func TestCodexCLIStatus_ModelUnsetWarns(t *testing.T) {
	bin := fakeLoginCLI(t, 0, "Logged in using ChatGPT")
	st := CodexCLIStatus(context.Background(), CodexCLIConfig{Bin: bin})
	if !st.Ready {
		t.Error("model 未設定でもログイン済みなら ready=true（注意喚起のみ）")
	}
	if !strings.Contains(st.Guidance, "CODEX_CLI_MODEL") {
		t.Errorf("guidance に CODEX_CLI_MODEL の注意が無い: %q", st.Guidance)
	}
}
