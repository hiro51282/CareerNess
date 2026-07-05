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

1. **codex login への動線（オンボーディング）が無い**
   実 AI を使うにはユーザーが自分で `codex login` 済みである必要があるが、アプリ内に案内・導線が無い。
   未認証時のエラー表示や「実 AI を使うには codex login が必要」のガイドが未整備。

2. **配布（distribution）**
   配布可能な実行形式に codex CLI を同梱できるか、できたとしてユーザーへ重いインストールを強要せずに
   済むかが未検討。→ **Desktop Host（Electron/Tauri/Wails）判断と連動**（`ai-foundation-direction.md`）。

3. **チャットが抽出専用で自由対話ができない**
   現状の message 経路は毎ターンを fact 抽出として扱うため、非 fact の発言
   （例：「他にどんな情報を書けばいいですか？」）で抽出 0 件となり `no facts extracted` エラーになる。
   会話としての往復・聞き返しができない。
   - 方向性（未決定・要設計）：(a) 0 件時をエラーにせずグレースフルな会話返信にする、
     (b) 「会話 vs 抽出」を判断するオーケストレーション層を設ける（Workspace Agent とは別責務）。
   - clarification questions（現状 source_detail に埋め込み）を一級の対話に昇格する検討とも関連。

## 関連
- `docs/implementation/ai/ai-foundation-direction.md`（正式経路・credential 非管理・Desktop First・Runtime 非抽象化）
- `AGENTS.md`（既知の未解決事項・Desktop Host 選定）
