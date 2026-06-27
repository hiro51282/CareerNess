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
// 選択ルール（Task4 Phase A）:
//   - EXTRACTION_PROVIDER 未設定 or "mock" → MockExtractionProvider（既定）
//   - "codex" → CodexExtractionProvider。OPENAI_API_KEY 欠落は明示エラー
//     （mock へのサイレントフォールバックはしない）
//   - それ以外 → 明示エラー
//
// credential（OPENAI_API_KEY）はここで env から読み、構築した provider に注入するのみ。
// 永続化せず、session / workspace とも混在させない。呼び出し側（handler）は本関数を
// リクエスト毎に呼ぶことで provider をリクエストスコープで構築する。これにより将来の
// マルチユーザー化（credential を env からリクエスト由来へ差し替える）を阻害しない。
func NewProviderFromEnv() (ExtractionProvider, error) {
	switch os.Getenv("EXTRACTION_PROVIDER") {
	case "", "mock":
		return NewMockExtractionProvider(), nil
	case "codex":
		return NewCodexExtractionProvider(CodexConfig{
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			BaseURL: os.Getenv("OPENAI_BASE_URL"), // 任意。空なら provider 既定
			Model:   os.Getenv("OPENAI_MODEL"),    // 任意。空なら provider 既定
		})
	default:
		return nil, fmt.Errorf("unknown EXTRACTION_PROVIDER %q (expected \"mock\" or \"codex\")", os.Getenv("EXTRACTION_PROVIDER"))
	}
}
