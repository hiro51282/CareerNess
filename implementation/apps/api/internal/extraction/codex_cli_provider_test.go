package extraction

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// fakeCodex は codex CLI を模した bash スクリプトを生成し、そのパスを返す。
// 引数から -o <path> を取り出し、mode に応じて payload をそこへ書く / 失敗する / 遅延する。
// 実 codex・ネットワーク不要で CI 安定に検証する。
func fakeCodex(t *testing.T, mode, payload string) string {
	t.Helper()
	dir := t.TempDir()
	payloadPath := filepath.Join(dir, "payload.txt")
	if err := os.WriteFile(payloadPath, []byte(payload), 0644); err != nil {
		t.Fatal(err)
	}
	scriptPath := filepath.Join(dir, "fakecodex.sh")
	script := fmt.Sprintf(`#!/usr/bin/env bash
out=""
while [ $# -gt 0 ]; do
  if [ "$1" = "-o" ]; then out="$2"; shift 2; continue; fi
  shift
done
case "%s" in
  ok)   cat %q > "$out" ;;
  fail) echo "boom" >&2; exit 3 ;;
  slow) sleep 5; cat %q > "$out" ;;
esac
`, mode, payloadPath, payloadPath)
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	return scriptPath
}

func TestCodexCLI_Happy(t *testing.T) {
	bin := fakeCodex(t, "ok", validExtractionJSON)
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: bin})

	res, err := p.ExtractFacts(context.Background(), "2022年にABCで決済基盤をGoへ移行した")
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}
	if len(res.ExtractedFacts) != 1 {
		t.Fatalf("facts = %d, want 1", len(res.ExtractedFacts))
	}
	if res.ExtractedFacts[0].FactIDHint != "payment-platform" {
		t.Errorf("fact_id_hint = %q", res.ExtractedFacts[0].FactIDHint)
	}
	// 会話返信 reply も契約の一部としてパースされること（B: 自由対話）
	if res.Reply == "" {
		t.Error("reply がパースされていない")
	}
	if p.Name() != "codex-cli" {
		t.Errorf("Name = %q, want codex-cli", p.Name())
	}
}

func TestCodexCLI_StripsFence(t *testing.T) {
	bin := fakeCodex(t, "ok", "```json\n"+validExtractionJSON+"\n```")
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: bin})

	res, err := p.ExtractFacts(context.Background(), "x")
	if err != nil {
		t.Fatalf("フェンス付き出力の解析に失敗: %v", err)
	}
	if len(res.ExtractedFacts) != 1 {
		t.Fatalf("facts = %d, want 1", len(res.ExtractedFacts))
	}
}

func TestCodexCLI_MalformedJSON(t *testing.T) {
	bin := fakeCodex(t, "ok", "これは JSON ではありません")
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: bin})

	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("不正な出力はエラーになるべき")
	}
}

func TestCodexCLI_NonZeroExit(t *testing.T) {
	bin := fakeCodex(t, "fail", "")
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: bin})

	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("非ゼロ終了はエラーになるべき")
	}
}

func TestCodexCLI_Timeout(t *testing.T) {
	bin := fakeCodex(t, "slow", validExtractionJSON) // sleep 3s
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: bin, Timeout: 100 * time.Millisecond})

	start := time.Now()
	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("timeout はエラーになるべき")
	}
	// timeout(100ms) + WaitDelay(2s) 程度で返るはず。fake の sleep 5s より十分早い。
	if elapsed := time.Since(start); elapsed > 4*time.Second {
		t.Errorf("timeout が効いていない（%v 経過）", elapsed)
	}
}

func TestCodexCLI_MissingBinary(t *testing.T) {
	p := NewCodexCLIProvider(CodexCLIConfig{Bin: "/nonexistent/codex-xyz"})
	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("バイナリ不在はエラーになるべき")
	}
}
