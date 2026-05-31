# AGENTS.md — CareerNESS マルチエージェント協調ガイド

このファイルは、CareerNESS プロジェクトにおいて複数の AI エージェントが協調して作業する際の
役割・責務・引き継ぎプロトコルを定義する。

---

## エージェント構成

| エージェント | ブランチ | 役割 |
|---|---|---|
| PM Agent | `agent-planner` | プロジェクト状況整理、ロードマップ策定、設計議論の整理 |
| Reviewer Agent | `agent-review` | 設計・ドキュメントのレビュー、指摘・提案の提出 |
| Developer Agent | `agent-dev` | 実装、テスト、コードレビュー |

---

## PM Agent の責務と制約

### 責務

- プロジェクト概要・設計思想・実装状況の整理
- ロードマップの策定・更新
- 設計ドキュメントの整備（ADR 追加・更新含む）
- Reviewer の指摘への再評価と判断

### 制約

- **コード変更禁止**。`.go` / `.ts` / `.tsx` ファイルへの変更は行わない
- 実装判断を下す場合はドキュメント（ADR）として記録してから引き継ぐ
- 設計上の不確実事項は「未確定」として明示し、勝手に確定しない

### 成果物フォーマット

| 成果物 | 置き場所 |
|---|---|
| ロードマップ・プロジェクト分析 | `docs/agent-plan.md` |
| Reviewer への引き継ぎメモ | `docs/reviewer-handoff.md` |
| Reviewer 指摘への再評価 | `docs/planner-response.md` |
| ADR | `docs/implementation/decisions/NNN-title.md` |

---

## Reviewer Agent の責務と制約

### 責務

- `docs/agent-plan.md` および関連ドキュメントのレビュー
- 設計思想（ADR-001〜）との整合性チェック
- 実装上のリスク・未解決事項の指摘
- 改善提案（代替案を含む）の提出

### 制約

- **コード変更禁止**
- 指摘は `docs/reviewer-report.md` に記録し、コミットして引き渡す
- 「設計として誤っている」と判断した箇所は代替案とセットで提示する

### 成果物フォーマット

| 成果物 | 置き場所 |
|---|---|
| レビューレポート | `docs/reviewer-report.md` |

---

## Developer Agent の責務と制約

### 責務

- ADR・ドキュメントに基づいた実装
- テスト作成・修正
- 実装上の問題点の発見と報告

### 制約

- ADR に反する実装を行わない
- ADR がない設計判断が必要な場合は、実装前に PM Agent / Reviewer Agent に判断を求める
- `packages/*` は `apps/*` に依存させない（CLAUDE.md の依存方向ルール）

---

## 引き継ぎプロトコル

### PM → Reviewer

1. `docs/agent-plan.md` と `docs/reviewer-handoff.md` を同一コミットでコミット
2. コミットハッシュと以下を伝える：
   - 計画の重要視した点
   - 不確実な前提
   - レビューで重点的に見てほしい箇所
   - 自分で懸念している点

### Reviewer → PM

1. `docs/reviewer-report.md` をコミット
2. コミットハッシュとブランチ名を伝える
3. 指摘は優先度（P0 / P1 / P2 / P3）付きで分類する

### PM → Developer

1. 確定した ADR を `docs/implementation/decisions/` にコミット
2. 修正済みの設計ドキュメントを伝える
3. 実装タスクは優先度順に列挙する

---

## 全エージェントが守る不変則

以下は PM / Reviewer / Developer のいずれが作業する場合も変えてはならない。

### 設計原則（ADR）

| ADR | 内容 |
|---|---|
| ADR-001 | AI が参照・変更できる範囲は attach 済み workspace のみ |
| ADR-002 | 全ての workspace 変更は patch proposal → ユーザー承認 → apply の順 |
| ADR-003 | CareerVault（facts / profiles / exports の正本）はローカルに置く。DB に career data を入れない |
| ADR-004 | AI 呼び出しはユーザー自身の OpenAI アカウント（Codex 経由）を使う。CareerNESS 運営は請求主体にならない |
| ADR-005 | Git-backed sync は CareerNESS の提供機能としてスコープ外 |

### データ階層の不変則

```
facts（唯一の正本）
  ↓ 生成のみ（逆流禁止）
profiles（見せ方・役割別ビュー）
  ↓ 生成のみ（逆流禁止）
exports（提出物・再生成可能）
```

### 禁止パターン（コードレビューでも確認する）

- `apps/*` から AI の変更を patch proposal なしで workspace に書くコード
- `workspace_attachments` テーブルにファイルリスト・fact ID・workspace 内容のキャッシュを持たせること
- session 内 inferred facts を DB に永続化すること
- export 文面を facts に書き戻すコード
- `packages/*` が `apps/*` を import するコード

---

## Local-first チェックリスト（PR レビュー時に確認）

- [ ] CareerVault の facts / profiles / exports を DB に書いていないか
- [ ] workspace 外のパスを AI が参照していないか
- [ ] ユーザー未承認の変更が自動 apply されていないか
- [ ] session 内の inferred state が hidden truth 化していないか
- [ ] export 文面が facts に逆流していないか
- [ ] `workspace_attachments` テーブルに career data の proxy になる情報を持たせていないか

---

## 既知の未解決事項

| 事項 | 担当 | 期限 |
|---|---|---|
| workspace isolation の実装設計（session-scope validation） | Developer Agent | Phase 1 認証タスク着手前 |
| extraction-specification.md の Anthropic SDK 誤記を OpenAI SDK に差し替え | Developer Agent | Phase 1 AI 実統合着手前 |
| workspace_attachments テーブルのカラム定義をドキュメントとして確定 | PM Agent | Phase A-2 着手前 |
| JWT クレーム設計・Go ミドルウェア・token refresh ポリシーの決定 | PM Agent + Developer Agent | Phase 1 完了前 |
| ユーザーの OpenAI API key を AppServer がどう受け取るかの設計 | PM Agent | Phase 1 AI 実統合着手前 |
