package extraction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CodexCLIProvider は credential-free の実 AI 正式経路（ai-foundation-direction.md）。
// ユーザーが独立に認証済みの codex CLI を非対話（exec）で実行し、構造化 JSON を得る。
// credential は codex CLI 側に閉じ、CareerNESS は鍵に一切触れない。
//
// CLI 実行は runCLI に隔離し、prompt 生成・ExtractedFactResult・JSON schema を知らない
// 疎結合を保つ。将来 Claude Code / Gemini CLI 等を足す際は runCLI 相当を Runtime 層へ
// 昇格させれば足りる構造にする（今は interface 化しない = YAGNI）。
type CodexCLIProvider struct {
	bin     string
	model   string
	timeout time.Duration
}

const (
	defaultCodexCLIBin     = "codex"
	defaultCodexCLITimeout = 60 * time.Second
)

// CodexCLIConfig は CodexCLIProvider の構築パラメータ。
type CodexCLIConfig struct {
	Bin     string        // 任意。既定 "codex"
	Model   string        // 任意。codex exec -m に渡す
	Timeout time.Duration // 任意。既定 60s
}

// NewCodexCLIProvider は config から provider を構築する。
func NewCodexCLIProvider(cfg CodexCLIConfig) *CodexCLIProvider {
	bin := cfg.Bin
	if strings.TrimSpace(bin) == "" {
		bin = defaultCodexCLIBin
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultCodexCLITimeout
	}
	return &CodexCLIProvider{bin: bin, model: cfg.Model, timeout: timeout}
}

// Name は provider 名を返す。
func (c *CodexCLIProvider) Name() string { return "codex-cli" }

// ExtractFacts は会話を codex CLI へ渡し、ExtractedFactResult を返す。
func (c *CodexCLIProvider) ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error) {
	prompt := getSystemPrompt() + "\n\n" + getUserPrompt(conversation)

	output, err := c.runCLI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	content := stripCodeFence(output)
	var result ExtractedFactResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("codex-cli: parse extraction JSON: %w", err)
	}
	return &result, nil
}

// runCLI は codex CLI を非対話実行し、モデルの最終メッセージ（テキスト）を返す。
//
// prompt は stdin で渡し（"-"）、`--output-last-message <tmp>` で最終メッセージのみを
// ファイル取得することで、エージェントの途中出力（stdout の chatter）を避ける。
// 抽出固有のこと（prompt 文言・ExtractedFactResult・schema）は一切知らない疎結合の seam。
func (c *CodexCLIProvider) runCLI(ctx context.Context, input string) (string, error) {
	tmp, err := os.CreateTemp("", "codex-out-*.txt")
	if err != nil {
		return "", fmt.Errorf("codex-cli: create temp: %w", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(tmpPath)

	runCtx := ctx
	if c.timeout > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	args := []string{"exec", "--skip-git-repo-check", "-o", tmpPath}
	if c.model != "" {
		args = append(args, "-m", c.model)
	}
	// プロンプトは stdin から受け取らせる。
	args = append(args, "-")

	cmd := exec.CommandContext(runCtx, c.bin, args...)
	cmd.Stdin = strings.NewReader(input)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// ctx キャンセルで親を kill しても、子孫が I/O パイプを握っていると Wait が
	// ブロックし得る。WaitDelay で kill 後の I/O 待ちを上限化してハングを防ぐ。
	cmd.WaitDelay = 2 * time.Second

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("codex-cli: exec failed: %w: %s", err, truncate(stderr.String(), 300))
	}

	out, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("codex-cli: read output: %w", err)
	}
	return string(out), nil
}
