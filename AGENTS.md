# AGENTS.md — CareerNESS Maker/Checker 協調ガイド

このファイルは、CareerNESS で AI エージェント（および人間）が **Maker/Checker モデル**で協調するための
役割・検証ゲート・引き継ぎ・不変則を定義する。

現状フェーズ: **CareerNESS Editor（Career Vault の編集ツール）MVP のコアループ完成**
（チャット → Fact 抽出 → レビュー → Vault 反映 → Facts 閲覧 → 確定）。実 AI（Codex CLI）接続は今後。
AI 基盤の方向性は `docs/implementation/ai/ai-foundation-direction.md` を参照。

---

## Loop Engineering / Maker–Checker

開発は「小さな増分 → 自動/人手で検証 → 緑ならマージ」の閉じたループで回す。

```
プラン提出 → ユーザー承認 → 実装(Maker) → PR
  → 自動 Checker（CI: api/web, CodeQL）→ 人手/設計 Checker（レビュー・不変則）
  → マージ（分類に応じ Claude 自己 or ユーザー承認）→ master 同期・整理
```

### Checker の構成

| レイヤ | 担い手 | 役割 |
|---|---|---|
| **Maker** | Developer（実装）＝現状 Claude Code | 実装・テスト作成・PR 作成 |
| **自動 Checker** | **CI（api / web）＋ CodeQL** | build/vet/test・型・ビルド・security を機械強制（merge ゲート） |
| **人手/設計 Checker** | Reviewer・PM＝現状ユーザー＋Claude | 設計・ADR 整合・Local-first チェックリスト・仕様判断 |

CI 定義は `.github/workflows/ci.yml`、PR チェックリストは `.github/pull_request_template.md`。

### DoD（完了条件）

- Go: `go build ./... && go vet ./... && go test ./...` 緑（`implementation/apps/api`）
- Web: `pnpm --filter web run test` と `pnpm build:web`（`tsc` 込み）緑（`implementation`）
- CI（api / web / CodeQL）緑
- Local-first / 設計不変則チェック通過（後述）

### マージ権限の分類

CI 緑・レビュー済みの PR は **Claude Code が自己マージ**する。次のみ**ユーザー承認**を要す：

1. **重大な変更**（広範囲・不可逆・高影響）
2. **セキュリティ要因**（workspace 境界・認証・権限・API キー/鍵・公開範囲・CodeQL 指摘）
3. **仕様変更に関わる判断**（設計・ADR・データ契約・プロダクト方針）
4. **その他**、ユーザーの判断/承認が必要と Claude が判断した事項

迷ったら例外（要承認）側に倒す。自己マージ時も分類を一言添えて透明化する。

---

## エージェント構成（論理的役割）

| エージェント | ブランチ | 役割 | Maker/Checker |
|---|---|---|---|
| PM Agent | `agent-planner` | 状況整理、ロードマップ、設計/ADR 整備 | 設計 Checker |
| Reviewer Agent | `agent-review` | 設計・ドキュメントレビュー、指摘・提案 | 設計 Checker |
| Developer Agent | `agent-dev` | 実装、テスト、コードレビュー | Maker |

> **現状は Claude Code が Maker、ユーザー＋Claude が人手 Checker、CI/CodeQL が自動 Checker** を担う。
> 上表の3エージェント×専用ブランチは**将来のマルチエージェント運用のための論理定義**であり、
> 現行の必須運用ではない（下記「引き継ぎプロトコル」参照）。

---

## PM Agent の責務と制約

### 責務

- プロジェクト概要・設計思想・実装状況の整理
- ロードマップの策定・更新
- 設計ドキュメントの整備（ADR 追加・更新含む。**ADR-006 および `ai-foundation-direction.md` を含む**）
- Reviewer の指摘への再評価と判断

### 制約

- **コード変更禁止**。`.go` / `.ts` / `.tsx` ファイルへの変更は行わない
- 実装判断を下す場合はドキュメント（ADR）として記録してから引き継ぐ
- 設計上の不確実事項は「未確定」として明示し、勝手に確定しない

---

## Reviewer Agent の責務と制約

### 責務

- 設計ドキュメント・ADR との整合性チェック（ADR-001〜006）
- 実装上のリスク・未解決事項の指摘
- 改善提案（代替案を含む）の提出

### 制約

- **コード変更禁止**
- 「設計として誤っている」と判断した箇所は代替案とセットで提示する
- 指摘は優先度（P0 / P1 / P2 / P3）付きで分類する

---

## Developer Agent の責務と制約

### 責務

- ADR・ドキュメントに基づいた実装
- テスト作成・修正、DoD の充足
- 実装上の問題点の発見と報告

### 制約

- ADR に反する実装を行わない
- ADR がない設計判断が必要な場合は、実装前に判断を求める（プラン提出 → ユーザー承認）
- `packages/*` は `apps/*` に依存させない（CLAUDE.md の依存方向ルール）

---

## 引き継ぎプロトコル

### 現行フロー（既定）

単一の Maker（Claude Code）＋人手 Checker（ユーザー）＋自動 Checker（CI/CodeQL）で回す：

```
プラン提出 → ユーザー承認 → 実装 → PR 作成 → CI/CodeQL（自動 Checker）
  → レビュー（不変則・設計） → マージ（分類に応じ Claude 自己 or ユーザー承認）
  → master 同期・ブランチ整理
```

### マルチエージェント handoff（任意 / legacy）

将来 PM / Reviewer / Developer を別エージェントで分離する場合の doc ベース引き継ぎ。
**現行の単一エージェント運用では使わない**（PR + CI + レビューが主）。

- PM → Reviewer: `docs/agent-plan.md` ＋ `docs/reviewer-handoff.md` を同一コミット。重視点・不確実な前提・重点レビュー箇所を伝える
- Reviewer → PM: `docs/reviewer-report.md` をコミット。指摘は P0〜P3 で分類
- PM → Developer: 確定 ADR を `docs/implementation/decisions/` にコミットし、実装タスクを優先度順に列挙

---

## 全エージェントが守る不変則

以下は誰が作業する場合も変えてはならない。

### 設計原則（ADR）

| ADR | 内容 |
|---|---|
| ADR-001 | AI が参照・変更できる範囲は attach 済み workspace のみ |
| ADR-002 | 全ての workspace 変更は patch proposal → ユーザー承認 → apply の順 |
| ADR-003 | CareerVault（facts / profiles / exports の正本）はローカルに置く。DB に career data を入れない |
| ADR-004 | AI 呼び出しはユーザー自身の OpenAI アカウント（Codex 経由）。CareerNESS 運営は請求主体にならない。**注: credential 非管理＋Codex CLI 正式経路へ転換予定（`ai-foundation-direction.md`、ADR 化候補）** |
| ADR-005 | Git-backed sync は CareerNESS の提供機能としてスコープ外 |
| ADR-006 | AI の workspace 変更は attach された root 配下に封じ込める（session 束縛 + `workspace.ResolveWithin` で強制）|

**AI 基盤のプロダクト原則**（`ai-foundation-direction.md`、ADR 化候補）:

- CareerNESS は AI provider の credential を管理しない（信頼境界）
- Structured Extraction の正式経路は **Codex CLI**（HTTP provider は deprecated・休眠保持）
- Desktop First は現時点のプロダクト方針（Desktop Host の実装技術は未決定・別判断）
- Runtime 抽象化はしない（YAGNI）。Workspace Agent は Provider の置き換えでなく別責務として将来共存

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
- **CareerNESS が AI provider の credential（API キー / トークン）を保持・保存する、または API キー入力を基本 UX 化するコード**（`ai-foundation-direction.md` のプロダクト原則）

---

## Local-first チェックリスト（PR レビュー時に確認）

> `.github/pull_request_template.md` に同期。PR 本文でも確認する。

- [ ] CareerVault の facts / profiles / exports を DB に書いていないか
- [ ] workspace 外のパスを AI が参照していないか
- [ ] ユーザー未承認の変更が自動 apply されていないか
- [ ] session 内の inferred state が hidden truth 化していないか
- [ ] export 文面が facts に逆流していないか
- [ ] `workspace_attachments` テーブルに career data の proxy になる情報を持たせていないか
- [ ] **CareerNESS が AI provider の credential を保持していないか**

---

## 既知の未解決事項

| 事項 | 状態 / 引き継ぎ先 |
|---|---|
| workspace isolation の実装設計（session-scope validation） | **完了**（Task2 / ADR-006） |
| ユーザーの OpenAI 鍵を AppServer がどう受け取るかの設計 | **credential 非管理へ転換**（`ai-foundation-direction.md`、ADR 化候補） |
| extraction-specification.md の AI provider 記述の訂正 | **一部完了**（reply / 履歴の契約・prompt は更新済み）。§4 の旧実装例（Anthropic SDK）は歴史的記述として残存（バナーで注記） |
| workspace_attachments テーブルのカラム定義の確定 | Task3（マルチユーザー / 認証）へ繰り延べ |
| JWT クレーム設計・Go ミドルウェア・token refresh ポリシー | Task3 へ繰り延べ |
| credential 非管理・Desktop First の ADR 化 | 近く起票（ADR-004 改訂 / supersede を含む） |
| Desktop Host の技術選定（Electron / Tauri / Wails） | Editor MVP 後の別判断 |
| Codex CLI provider の実装（実 AI・正式経路） | **完了**（PR #19。実 AI で動作確認済み。有効化手順は `docs/implementation/ai/codex-cli-integration.md`） |
| codex login への動線（オンボーディング） | 次期候補（`codex-cli-integration.md` 課題#1） |
