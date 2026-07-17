# Getting Started — CareerNESS を最新で動かすまで

最新の master を取得してから、アプリを起動するまでの基本手順。
開発コマンドの一覧は `CLAUDE.md`、実 AI の詳細は `docs/implementation/ai/codex-cli-integration.md` を参照。

## 前提ツール

| ツール | バージョン | 備考 |
|---|---|---|
| Go | 1.23+ | backend（`implementation/apps/api`） |
| Node.js | 20+ | フロント / Electron |
| pnpm | 9+ | `~/.local/bin/pnpm`（このリポジトリの標準） |
| codex CLI | 任意 | **実 AI を使う場合のみ**。`codex login` 済みであること |

## 初回セットアップ

```bash
git clone <repo> && cd CareerNess/implementation
~/.local/bin/pnpm install
```

## 起動方法（3通り）

### A. ブラウザ開発モード（UI 開発向け・最速）

```bash
# 端末1: Go API サーバ（:8080）
cd implementation/apps/api && go run ./cmd/server

# 端末2: Vite dev サーバ（:5173）
cd implementation && ~/.local/bin/pnpm dev:web
```
Chrome / Edge で `http://localhost:5173`。ワークスペースの読み書きは
File System Access API（ブラウザ経路）で行われる。

### B. デスクトップ開発モード（Electron・Go 経路）

```bash
cd implementation
~/.local/bin/pnpm build:web                                  # フロントをビルド
(cd apps/api && go build -o bin/server ./cmd/server)         # Go バイナリ
cd apps/desktop && ~/.local/bin/pnpm run build && ~/.local/bin/pnpm start
```
ウィンドウが開き、Go サーバは子プロセスとして自動起動する（手動起動不要）。
attach はネイティブのフォルダ選択、読み書きは Go 経路（ADR-006/008）。

### C. パッケージ版（配布物・codex 同梱）

```bash
cd implementation && ~/.local/bin/pnpm build:web
cd apps/desktop && ~/.local/bin/pnpm run dist:linux          # 2〜3 分
./release/CareerNESS-0.0.1.AppImage
```
`dist:linux` は Go ビルド・リソース staging（codex native 同梱を含む）まで全て行う。

## 実 AI（Codex CLI）を使う

既定は Mock。実 AI は env で有効化する（**モデル指定は実質必須**。ChatGPT アカウントは `gpt-5.4` 等）:

```bash
EXTRACTION_PROVIDER=codex-cli CODEX_CLI_MODEL=gpt-5.4 <上記いずれかの起動コマンド>
```

- 事前に `codex login` が必要（チャット上部の AI バッジ／案内バナーで状態を確認できる）
- パッケージ版（C）は codex を同梱しているため、実 AI に必要なのは env と login のみ
- 詳細・トラブルシュート: `docs/implementation/ai/codex-cli-integration.md`

## 最新化と再ビルドの対応表

`git pull` した後、**何が変わったかで再ビルド対象が決まる**:

| 変更されたもの | ブラウザ dev (A) | デスクトップ dev (B) | パッケージ (C) |
|---|---|---|---|
| フロント（apps/web） | 不要（HMR） | `pnpm build:web` | `dist:linux` 再実行 |
| Go（apps/api） | サーバ再起動 | `go build -o bin/server ./cmd/server` | `dist:linux` 再実行 |
| Electron main（apps/desktop） | — | `pnpm run build`（tsc） | `dist:linux` 再実行 |

迷ったら **C は `dist:linux` を再実行すれば常に全部入り**（フロントだけ事前に `build:web` が必要）。

## 動作の一巡（何ができるか）

1. 「CareerVault を開く」でフォルダを attach
2. チャットでキャリアを話す → AI が会話しつつ fact 候補を提案（既存 Vault の facts も参照する）
3. 質問チップに答えると同じ fact が更新提案される
4. レビューで承認 → `facts/*.yaml` に書き込み
5. Facts タブで閲覧・「確定にする」で `proposed` → `confirmed`

## 検証（開発時）

PR 前のローカル検証は `CLAUDE.md` の「検証ゲート / DoD」を参照
（Go: build/vet/test、Web: vitest/build、CI が自動 Checker）。
