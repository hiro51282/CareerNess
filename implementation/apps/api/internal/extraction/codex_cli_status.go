package extraction

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// CLIStatus は codex CLI の利用可否診断の結果（オンボーディング用）。
type CLIStatus struct {
	Ready    bool   `json:"ready"`
	Detail   string `json:"detail"`
	Guidance string `json:"guidance,omitempty"`
}

// CodexCLIStatus は codex CLI が実 AI 経路として使える状態かを安価に診断する
//（`codex login status` を使うためトークン消費なし）。
//
// ExtractionProvider の責務（抽出）には含めず、独立した診断関数として提供する
//（provider interface を広げない方針。ai-foundation-direction.md）。
// 判定は 3 段: バイナリ存在 → ログイン状態 → モデル設定（未設定は ready のまま注意喚起）。
func CodexCLIStatus(ctx context.Context, cfg CodexCLIConfig) CLIStatus {
	bin := cfg.Bin
	if strings.TrimSpace(bin) == "" {
		bin = defaultCodexCLIBin
	}

	if _, err := exec.LookPath(bin); err != nil {
		return CLIStatus{
			Ready:    false,
			Detail:   "codex CLI が見つかりません",
			Guidance: "codex CLI のインストールが必要です。インストール後に「再確認」を押してください。",
		}
	}

	statusCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(statusCtx, bin, "login", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.WaitDelay = 2 * time.Second
	if err := cmd.Run(); err != nil {
		return CLIStatus{
			Ready:    false,
			Detail:   strings.TrimSpace(out.String()),
			Guidance: "codex にログインしていません。ターミナルで `codex login` を実行し、完了後に「再確認」を押してください。",
		}
	}

	st := CLIStatus{Ready: true, Detail: strings.TrimSpace(out.String())}
	if strings.TrimSpace(cfg.Model) == "" {
		// ChatGPT アカウントでは codex の既定モデルが使えないため、未設定は注意喚起する
		//（API キー認証では既定でも動き得るため ready は維持）。codex-cli-integration.md 参照。
		st.Guidance = "CODEX_CLI_MODEL が未設定です。ChatGPT アカウントでは既定モデルが使えないため、利用可能なモデル（例: gpt-5.4）の指定を推奨します。"
	}
	return st
}
