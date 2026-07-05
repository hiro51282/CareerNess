package extraction

import (
	"fmt"
	"os"
)

// NewProviderFromEnv は環境変数に基づいて extraction provider を選択・構築する。
//
// 責務は「どの provider を返すか」に限定する。timeout / retry / logging などの
// 横断的関心事はここに持ち込まない（provider 自身および将来の Phase で扱う）。
//
// TODO(Phase B): Phase B 以降は credential を env 固定ではなく request / config から
// 受け取る方向になる想定。そのタイミングで本関数を NewProvider(config) のような
// 名称・シグネチャへ整理するか再評価する（Phase A では env 固定で適切なため変更不要）。
//
// 選択ルール:
//   - EXTRACTION_PROVIDER 未設定 or "mock" → MockExtractionProvider（既定）
//   - "codex-cli" → CodexCLIProvider（**正式な実 AI 経路**。credential は codex CLI 側に
//     閉じ CareerNESS は鍵に触れない。ai-foundation-direction.md）
//   - "codex" → CodexExtractionProvider（HTTP/OpenAI 直叩き。**deprecated**：CareerNESS が
//     credential を管理する前提のため現行 MVP の正式経路から外している。休眠保持）
//   - それ以外 → 明示エラー（mock へのサイレントフォールバックはしない）
//
// 呼び出し側（handler）は本関数をリクエスト毎に呼ぶことで provider をリクエストスコープで
// 構築する。これにより将来のマルチユーザー化を阻害しない。
func NewProviderFromEnv() (ExtractionProvider, error) {
	switch os.Getenv("EXTRACTION_PROVIDER") {
	case "", "mock":
		return NewMockExtractionProvider(), nil
	case "codex-cli":
		// credential-free。bin/model は任意（空なら provider 既定 "codex"）。
		return NewCodexCLIProvider(CodexCLIConfig{
			Bin:   os.Getenv("CODEX_CLI_BIN"),
			Model: os.Getenv("CODEX_CLI_MODEL"),
		}), nil
	case "codex":
		// deprecated（HTTP）。credential 管理前提のため正式経路から外す。
		return NewCodexExtractionProvider(CodexConfig{
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			BaseURL: os.Getenv("OPENAI_BASE_URL"), // 任意。空なら provider 既定
			Model:   os.Getenv("OPENAI_MODEL"),    // 任意。空なら provider 既定
		})
	default:
		return nil, fmt.Errorf("unknown EXTRACTION_PROVIDER %q (expected \"mock\", \"codex-cli\", or \"codex\")", os.Getenv("EXTRACTION_PROVIDER"))
	}
}
