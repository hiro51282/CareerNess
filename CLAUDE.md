# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業するためのガイダンスを提供します。

## 技術スタック

| 領域 | 技術 |
|---|---|
| フロントエンド | React + TypeScript + Vite（Next.js ではない） |
| バックエンド | Go 1.23.4（workspace runtime として扱う） |
| モノレポ | pnpm workspaces（Turborepo/Nx は使わない） |
| ワークスペース読み書き | File System Access API（ブラウザ側、Chrome/Edge のみ） |
| データ形式 | YAML |
| AI 統合 | Structured Extraction Provider（現状 Mock、正式経路は Codex CLI）。`docs/implementation/ai/ai-foundation-direction.md` 参照 |

pnpm は `~/.local/bin/pnpm` に入っている。

## 開発コマンド

```bash
# Go API サーバー起動（:8080）
cd implementation/apps/api && go run ./cmd/server

# React dev サーバー起動（:5173）
cd implementation && ~/.local/bin/pnpm dev:web

# React ビルド
cd implementation && ~/.local/bin/pnpm build:web

# Go ビルド確認
cd implementation/apps/api && go build ./...

# Go テスト
cd implementation/apps/api && go test ./...

# TypeScript 型チェック
cd implementation/apps/web && ~/.local/bin/pnpm exec tsc --noEmit
```

vite.config.ts に `/api` → `localhost:8080` のプロキシ設定済みなので、フロントエンドから `/api/v1/...` でバックエンドに接続できる。

### 検証ゲート / DoD

PR は自動 Checker（CI）と以下ローカル検証が緑であること（Loop Engineering。詳細は `AGENTS.md`）：

- Go: `go build ./...` / `go vet ./...` / `go test ./...`（`implementation/apps/api`）
- Web: `~/.local/bin/pnpm --filter web run test`（vitest）/ `~/.local/bin/pnpm build:web`（`tsc` 込み）
- CI: `.github/workflows/ci.yml`（`api` / `web`）＋ CodeQL。PR テンプレは `.github/pull_request_template.md`。

## プロジェクト状況

CareerNess は **CareerNESS Editor（Career Vault の編集ツール）MVP のコアループ完成**フェーズです。
チャット → Fact 抽出 → レビュー → Vault 反映 → Facts 閲覧 → 確定（`proposed` → `confirmed`）まで app 内で一巡できる。

- 実装は `apps/api`（Go: `internal/` の extraction / patch / workspace / session / handler / ai）と `apps/web`（React）にある（`packages/*` は現状スタブ）。
- Fact 抽出は現状 **Mock Provider**。実 AI の正式経路は **Codex CLI**（今後実装）。方向性は `docs/implementation/ai/ai-foundation-direction.md`。
- 開発プロセスは **Loop Engineering（Maker/Checker）**。`AGENTS.md` と CI（`.github/workflows/ci.yml`）を参照。

## リポジトリ構成

```
CareerNess/
├── docs/                        # 設計・仕様ドキュメント
│   ├── architecture/            # システム概要、責務境界
│   ├── workspace/               # CareerVault データモデルとレイアウト
│   ├── ai/                      # AI 動作仕様、プロンプト戦略
│   ├── implementation/          # 実装設計書（読む順序は後述）
│   └── ...
└── implementation/              # ソースコード（モノレポ。apps/ に実装あり、packages/ は現状スタブ）
    ├── apps/
    │   ├── web/                 # ブラウザ UI
    │   └── api/                 # AppServer — AI オーケストレーション層
    ├── packages/
    │   ├── workspace-core/      # ワークスペースアクセス、バリデーション、apply/rollback
    │   ├── patch-engine/        # セマンティックパッチモデルと差分ヘルパー
    │   ├── schema/              # Fact / Profile / Export のデータ契約
    │   └── shared/              # 最小限の共通ユーティリティのみ
    ├── infra/                   # Docker、Terraform、AWS デプロイ定義
    ├── scripts/                 # 開発補助スクリプト
    └── examples/                # サンプルワークスペースデータ
```

**依存方向のルール：** `apps/*` は `packages/*` に依存してよい。`packages/*` は `apps/*` に依存してはならない。`packages/shared` は薄く保つ — ゴミ箱にしない。

## コアアーキテクチャ

### 3層構成

```
Browser UI  →  Codex AppServer  →  CareerVault（ローカルワークスペース）
```

- **Browser**：会話・ワークスペース attach・パッチレビュー・閲覧の UI。一時的な UI state のみ保持。キャリアデータの正本は持たない。
- **AppServer** (`apps/api`)：認証セッション、AI オーケストレーション、パッチ提案/適用の仲介を担う。オーケストレーション層であり、キャリアの正本を所有しない。
- **CareerVault**：ユーザーが所有するローカルワークスペース（YAML / markdown ファイル）。唯一の正本 (source of truth)。

### パッチ提案モデル（ADR-002）

AI によるワークスペース変更はすべてパッチ提案を経由する。直接書き込みは禁止。

```
AI による抽出/生成
  → パッチ提案（YAML、人間・機械の両方が読める）
  → バリデーション（スキーマ、パス境界、ステータス遷移）
  → ユーザーレビュー（diff + 根拠 + 確信度）
  → 承認 / 却下 / 修正
  → ワークスペースへ適用
  → 履歴記録
```

AI は **提案のみ** 行い、単独で適用することはできない。`confirmed` ステータスはユーザー承認後にのみ設定される。

### ワークスペース境界（ADR-001）

AI の能力は明示的に attach されたワークスペースルートのみにスコープされる。attach 外のパス（ホームディレクトリ等）は読み取り・パッチ対象にならない。

### ローカルファースト（ADR-003）

CareerVault は facts / profiles / exports を人間可読な YAML ファイルとして保持する。AppServer が保持するのはセッションメタデータと一時的なオーケストレーション状態のみ。クラウドは認証と AI アクセスの補助に限定する。

## データモデル

以下の3層を混同してはならない：

| 層 | 役割 | 正本度 |
|---|---|---|
| **Fact** | 確認済みのキャリア事実（行動・判断・成果） | 最も正本に近い |
| **Profile** | 特定の audience/役職向けに facts を再構成したもの | 派生ビュー — 正本ではない |
| **Export** | 最終整形出力（履歴書、返信文など） | 提出用成果物 — 正本ではない |

補助概念：**Project**（facts をプロジェクト単位で束ねる）、**Tag**（横断的な分類ラベル）、**Role**（責務パターン、派生概念 — primary truth ではない）。

AI は fact を発明してはならない。不確実な項目は `proposed` または `inferred` のまま残す。`confirmed` への遷移にはユーザーの明示的な承認が必要。

## パッチフォーマット

パッチは YAML エンベロープ。必須フィールド：`patch_id`、`workspace_id`、`session_id`、`created_by`、`status`、`summary`、`operations` リスト。各 operation は `type`、`target`、`change`、`rationale`、`confidence` を持つ。fact に触れる operation はさらに `entity_id` と `fact_status_after` が必要。

MVP の operation タイプ：`create_file`、`update_file`、`delete_file`、`upsert_fact`、`mark_fact_status`、`replace_generated_profile`、`replace_generated_export`、`append_history_record`。

**1 patch = 1 セマンティック変更。** fact 追加とプロフィール再生成を1つのパッチにまとめない。

## パッケージ責務

### `packages/workspace-core`
担当：ワークスペースの attach/読み取り、YAML ファイル管理、スキーマバリデーション、patch apply、rollback、履歴。  
禁止：AI オーケストレーション、クラウド永続化、UI state、無制限ファイルスキャン。

### `packages/patch-engine`
担当：セマンティックパッチモデル、差分/バリデーションヘルパー、アトミック変更表現。  
禁止：ワークスペース正本の所有、クラウドセッション処理、UI レンダリング関連。

### `packages/schema`
担当：Fact / Profile / Export のデータ契約、ランタイム横断のバリデーションルール。  
禁止：AI プロンプトロジック、ワークスペース変更オーケストレーション。

### `apps/api`
担当：セッション管理、AI オーケストレーション、パッチ提案生成、ワークスペースゲートウェイ（読み取り専用リスト、apply 仲介）。  
禁止：facts / profiles / exports の恒久保存、ワークスペースオーナーとしての振る舞い。

## 設計上の必須制約

- **サイレント apply 禁止**：低リスクの変更（export 再生成、メタデータ更新）でもユーザーに見えない自動適用はしない。
- **Fact の捏造禁止**：AI は推論だけで `confirmed` の fact を追加してはならない。不明フィールドは空のまま提案する。
- **提案 ≠ 適用**：`generate_patch_proposal` と `apply_approved_patch` は常に別の操作。
- **ワークスペース外参照禁止**：AI ツールは attach 済みワークスペースルート外のパスを拒否する。
- **`packages/*` は `apps/*` に依存しない**：コアドメインロジックはエントリーポイントから切り離す。

## ドキュメントの読む順序

実装設計書は `docs/implementation/` を以下の順で読む：
1. `decisions/`（ADR-001、002、003）
2. `workspace/`（データモデル、ワークスペースレイアウト）
3. `ai/`（パッチフォーマット、承認フロー、ツールモデル）
4. `profile/`、`export/`
5. `backend/`、`frontend/`、`auth/`、`deploy/`
