# Codex CLI 統合（運用ノート）

> 現時点の状態の記録（2026-07-06）。正式経路は Codex CLI（`ai-foundation-direction.md`）。
> 実装は `internal/extraction/codex_cli_provider.go`、選択は `EXTRACTION_PROVIDER` env。

## 実 AI の有効化手順

### 前提
- `codex` CLI がインストール済みで、**`codex login` 認証済み**であること。
  credential は codex 側（`~/.codex`）に閉じ、**CareerNESS は鍵に触れない**（プロダクト原則）。

### 起動
```bash
# API サーバ（codex-cli provider）
cd implementation/apps/api
EXTRACTION_PROVIDER=codex-cli CODEX_CLI_MODEL=gpt-5.4 go run ./cmd/server   # :8080

# フロント（別端末）
cd implementation && ~/.local/bin/pnpm dev:web                              # :5173
# → Chrome で localhost:5173（ブラウザ dev では両サーバ必須）
```

### モデル指定（重要）
- **`CODEX_CLI_MODEL` は実質必須**。codex の既定 `gpt-5.3-codex` は **ChatGPT アカウント連携では非対応**
  （`*-codex` / reasoning 系 `o3`/`o4-mini` 等も同様に拒否される）。
- **ChatGPT アカウントでは素の `gpt-5.4` が使える**（実 smoke で確認済み）。プラン/アカウントにより
  使えるモデルは異なるため、各自の環境で通るモデルを指定する。

### 確認済み（2026-07-06 smoke）
- `/extract`（および chat の message 経路）で **実 codex 抽出が end-to-end 動作**。
  会話から experience fact を保守的に抽出（期間を捏造せず clarification questions を生成）。
- レイテンシ 約 31s（provider 既定 timeout 60s 内）。credential-free・config 駆動でコード変更ゼロ。

## 既知の課題 / バックログ（実 AI 利用で顕在化）

MVP のコアループ完成後に取り組む候補。実装判断はプラン→承認で行う。

1. **codex login への動線（オンボーディング）** → **検知＋案内は実装済み（2026-07-19）**
   `GET /api/v1/ai/status` が provider の利用可否を診断（バイナリ存在 → `codex login status` →
   モデル設定。トークン消費なし）。UI はハイブリッド方式＝起動時 1 回取得＋バナーの「再確認」
   ボタン＋送信エラー時の自動再取得で、バッジ（Mock / codex-cli ✅ / ⚠）と guidance を表示する。
   毎メッセージの事前チェックは行わない（外部プロセス起動のコストと冗長性のため）。
   **アプリ内から `codex login` を起動する動線は未実装**（ブラウザからは device code フローを
   中継できないため、Desktop Host（課題#2）判断とセットで将来対応）。

2. **配布（distribution）** → **Host 方針決定（2026-07-16）**：Desktop Host は **Electron**（ADR-008、
   thin-main 原則・PoC ゲート付き）。codex CLI の同梱可否／インストール誘導の具体設計は
   Electron PoC・配布タスクで検証する（未実装）。

3. ~~チャットが抽出専用で自由対話ができない~~ → **解決済み（2026-07-19, PR #21/#22/#23）**
   - PR #21: AI 出力契約に `reply`（会話返信）を追加し 0 件抽出を許容。非 fact の発言は会話として応答
     （`extraction-specification.md` §3/§5 が SSOT。`/extract` の 0 件→422 は不変）。
   - PR #22: 会話履歴（transcript・直近 10 ターン）を導入。既出 fact への詳細追加は同じ
     `fact_id_hint` を再利用 → 同一 fact_id の upsert で **聞き返しに答えると fact が育つ**。
   - PR #23: clarification questions を `clarifications[]` として応答に集約し、チャット UI に
     質問チップとして一級表示（タップで回答テンプレをプリフィル）。
   - 実 AI（gpt-5.4）で会話・抽出・enrich の全モードを smoke 済み。

## 関連
- `docs/implementation/ai/ai-foundation-direction.md`（正式経路・credential 非管理・Desktop First・Runtime 非抽象化）
- `AGENTS.md`（既知の未解決事項・Desktop Host 選定）
