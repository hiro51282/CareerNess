# Roadmap

CareerNess の roadmap は、MVP を中心に段階的に考える。最初から理想形を作るのではなく、設計思想を壊さずに最小限の流れを成立させることを優先する。

---

## 現状 (2026-05-31)

**Phase 0 完了、Phase 1 進行中。**

実装済み：
- モノレポ構成（pnpm workspaces + Go + React/Vite）
- Go API サーバー（`:8080`）
  - `POST /api/v1/conversations/message` — patch proposal 生成（mock AI）
  - `POST /api/v1/patches/validate` — patch バリデーション
  - `POST /api/v1/extract` — extraction pipeline（mock provider）
  - `POST /api/v1/apply-patch` — patch apply（workspace への YAML 書き込み）
- React Web アプリ（`:5173`）
  - ワークスペース attach（File System Access API）
  - サイドバー — attach 済みワークスペースのファイル一覧
  - チャット UI
  - パッチレビュー画面（diff 表示・承認・却下）
  - 承認後の YAML 書き込み（`facts/inbox.yaml`）
- patch モデル（`patch_id`、`operations`、`status`、`rationale`、`confidence`）
- patch バリデーション（パストラバーサル検出、operation タイプ制限）
- サンプルワークスペース（`implementation/examples/careervault/`）

未着手（Phase 1 残り）：
- AI の実統合（OpenAI Codex 経由。ADR-004 で確定）
- fact schema minimal fix（`action / decision / impact` を空フィールドとして追加。Phase 2 マイグレーション回避のため）
- fact 整形 UI（inbox の raw_text を構造化する）
- profile 生成
- export
- 認証（JWT + workspace isolation + ログイン画面）
- E2E 統合テスト（extract → patch → apply の full loop）

---

## Phase 0: Concept Fixing ✅

目的:

- Local-first の意味を明確にする
- workspace boundary を文書化する
- facts / profiles / exports の責務を固定する
- AI capability boundary を定義する

この phase では、実装よりも設計固定を優先する。

## Phase 1: MVP Workspace Flow 🔨 進行中

目的:

- ~~workspace attach~~ ✅
- ~~chat~~ ✅
- ~~patch proposal / review~~ ✅
- ~~YAML apply~~ ✅
- fact schema minimal fix（`action / decision / impact` を空フィールドとして追加）
- AI 実統合（OpenAI Codex 経由。ADR-004）
- fact 整形（inbox → 構造化 fact）
- profile 生成
- export
- 認証（JWT + per-user workspace isolation）
- E2E 統合テスト

この phase で成立すべきこと:

- attach された workspace だけを対象に AI が動く
- facts を review しながら保存できる
- profile を facts から生成できる
- export を派生物として扱える
- ユーザー認証済みセッションのみが API を操作できる
- session scope 外の workspace_path への apply-patch が拒否される

## Phase 1.5: packages 分離スプリント（Phase 1 完了後・Phase 2 着手前）

目的:

- `packages/schema` — YAMLFact / Patch 型を一元化（二重定義解消）
- `packages/patch-engine` — validator / applier を移行
- `packages/workspace-core` — file I/O / attach 管理を移行
- `apps/api` は orchestration + handler のみに絞る

この sprint を入れる理由: Feature A（マルチユーザー）で認証ミドルウェアと workspace isolation を追加する際、現在の mixed 構造では変更影響範囲が不明確になるため。

## Phase 2: Better Structure and Review

目的:

- fact extraction 精度の改善
- patch proposal の見やすさ改善
- taxonomy / role / tag の整理
- 未確定情報の扱い改善

ここでは AI 品質よりも、「間違いをどう止めるか」を改善する。

## Phase 3: Reuse and Iteration

目的:

- role 別 profile の再利用
- 既存 facts からの再生成
- export variation の拡充
- workspace の運用しやすさ向上

## Phase 4: Optional Cloud Enhancements

目的:

- 必要最小限の同期や補助機能の検討
- 運用監視の整備
- deploy の安定化

前提:

- 正本が cloud に移るわけではない
- Local-first と user-owned data を壊さない

## 明確に後回しにするもの

- enterprise HR platform 化
- local LLM mandatory 化
- GPU hosting
- huge cloud storage
- 大規模 infra 最適化

## スコープ外と確定したもの（ADR）

- **Git-backed sync**（ADR-005）: Local-first と緊張関係にあるため CareerNESS の提供機能から除外。ユーザーが自身で Git 管理するのは妨げない
- **Anthropic Claude API の直接使用**: AI provider は OpenAI Codex 経由（ADR-004）で確定。extraction-specification.md の Anthropic SDK 実装例は誤記
