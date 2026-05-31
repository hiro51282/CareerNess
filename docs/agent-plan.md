# CareerNESS Agent Plan

作成日: 2026-05-31  
作成者: PM エージェント (Claude Sonnet 4.6)  
対象ブランチ: agent-planner

---

## 1. プロジェクト概要整理

### 何を作っているのか

CareerNESS は **Local-first AI 支援キャリア構造化プラットフォーム**。  
転職活動で繰り返し発生する「同じ経歴を媒体ごとに書き直す」問題を、キャリア情報を構造化データとして管理することで解決する。

### 解決する問題

| 問題 | 現状 | CareerNESSが目指す状態 |
|---|---|---|
| 経歴を媒体ごとに書き直す | 職務経歴書・スカウト返信・面接メモが別々 | facts を一元管理し、用途別に生成 |
| 「見せ方」と「事実」が混在 | レジュメ文面がキャリアの正本になってしまう | facts / profiles / exports を明示分離 |
| AI に書かせると根拠が消える | AI生成文章から事実を逆引きできない | proposed → confirmed の承認フロー |
| データが SaaS に保管される | サービス停止でデータ消失リスク | ローカルワークスペース (CareerVault) が正本 |

### コアコンセプト

```
CareerVault（ローカルYAML）
  ├── facts/          ← 事実の正本
  ├── profiles/       ← 用途別の「見せ方」（派生）
  ├── exports/        ← 提出物（派生）
  └── projects/       ← facts をプロジェクト単位で束ねる
```

**AI は提案するだけ。適用はユーザーが承認してから。**

---

## 2. 設計思想整理

### 3大原則（ADR）

| ADR | 判断 | 理由 |
|---|---|---|
| ADR-001 Workspace Boundary | AI は attach 済みワークスペース外にアクセスしない | 能力境界を明確化し、所有権を守る |
| ADR-002 Patch-Oriented Editing | 全変更は patch proposal → ユーザー承認 → apply の順 | hallucinated fact の流入阻止、ロールバック可能性 |
| ADR-003 Local-First Storage | CareerVault は人間可読 YAML、正本はローカル | user-owned data の実現 |

### データ階層の不変則

```
facts（唯一の正本）
  ↓ 生成
profiles（見せ方・役割別ビュー）
  ↓ 生成
exports（提出物・再生成可能）
```

**この順序を逆流させない。** export 文面を facts に書き戻さない。profile を正本扱いしない。

### プラットフォーム DB と CareerVault の分離

| DB に置いてよいもの | DB に置いてはいけないもの |
|---|---|
| user account / auth session | canonical facts / profiles / exports |
| subscription / billing | workspace mirror / cache |
| usage control / rate limit | hidden career truth |

---

## 3. 現在の実装状況整理

### Phase 状況（2026-05-31 時点）

| Phase | 状態 |
|---|---|
| Phase 0: Concept Fixing | ✅ 完了 |
| Phase 1: MVP Workspace Flow | 🔨 進行中（約 60% 完了） |

### バックエンド（Go API: `implementation/apps/api/`）

| モジュール | 状態 | 備考 |
|---|---|---|
| `internal/extraction/` | ✅ 実装済み | MockExtractionProvider のみ（real AI 未接続） |
| `internal/patch/` | ✅ 実装済み | applier, validator, model |
| `internal/workspace/schema.go` | ✅ 実装済み | MVP スキーマバリデーション |
| `internal/handler/` | ✅ 実装済み | message, extract, apply-patch, patch endpoints |
| `internal/ai/propose.go` | 🟡 スタブ | real AI 統合未実装 |
| テスト | 🟡 部分的 | validator_test, normalizer_test は pass。applier_test は ready だが E2E なし |
| 認証 | ❌ 未着手 | API 保護なし |

#### 実装済み API エンドポイント

```
GET  /health
POST /api/v1/conversations/message   ← mock AI による patch proposal 生成
POST /api/v1/patches/validate        ← patch バリデーション
POST /api/v1/extract                 ← extraction pipeline（mock provider）
POST /api/v1/apply-patch             ← patch apply（workspace への YAML 書き込み）
```

### フロントエンド（React + Vite: `implementation/apps/web/`）

| コンポーネント | 状態 | 備考 |
|---|---|---|
| `WorkspaceAttach` | ✅ 実装済み | File System Access API（Chrome/Edge のみ） |
| `WorkspaceContext` | ✅ 実装済み | ファイル一覧、dirHandle 管理 |
| `ChatView` | ✅ 実装済み | メッセージ送受信、patch proposal 受取り |
| `PatchReview` | ✅ 実装済み | diff 表示、承認・却下、YAML 書き込み |
| fact 整形 UI | ❌ 未着手 | inbox.yaml の raw_text を構造化する画面 |
| profile 生成 UI | ❌ 未着手 | |
| export UI | ❌ 未着手 | |
| ログイン画面 | ❌ 未着手 | |

### パッケージ群（`implementation/packages/`）

| パッケージ | 状態 | 備考 |
|---|---|---|
| `workspace-core` | ❌ README のみ | 実装は apps/api 内に混在 |
| `patch-engine` | ❌ README のみ | 実装は apps/api 内に混在 |
| `schema` | ❌ README のみ | 実装は apps/api 内に混在 |
| `shared` | ❌ README のみ | |

### サンプルデータ

`implementation/examples/CareerVault/` — facts, profiles, exports, projects を含む最小ワークスペース。実装リファレンスとして利用可能。

---

## 4. 未実装機能整理

### Phase 1 残り（MVP 完成に必要）

| 機能 | 概要 | 工数概算 |
|---|---|---|
| **AI 実統合** | Codex AppServer 経由での real fact extraction。CodexExtractionProvider の実装 | L (1週間) |
| **fact 整形 UI** | inbox.yaml の raw_text をユーザーが構造化 fact に変換するインタラクション | M (2-3日) |
| **profile 生成** | confirmed facts → 役割別プロフィール生成（API + UI） | M (2-3日) |
| **export** | profiles → markdown/テキスト出力（API + UI） | M (2-3日) |
| **認証** | ログイン・セッション管理、API 保護 | L (1週間) |
| **E2E 統合テスト** | extraction → patch → apply の full loop テスト | S (1日) |

### Phase 2 以降（品質・利便性向上）

| 機能 | 概要 |
|---|---|
| セマンティック重複検出 | 同一 fact の重複提案を防ぐ |
| multi-stage extraction | 会話を重ねて facts を段階的に精緻化 |
| タイムライン UI | facts を時系列で可視化 |
| rollback 機能 | patch 適用前の状態に戻す |
| fact schema 高度化 | action / decision / impact / context / evidence の豊富なフィールド対応 |
| streaming AI 応答 | 長い AI 応答を stream 表示 |

---

## 5. 技術的負債整理

### 高優先度（MVP 前に対処すべき）

| 負債 | 詳細 | 対処方針 |
|---|---|---|
| **Patch struct の二重定義** | `internal/extraction/models.go` と `internal/patch/model.go` に別々の Patch/Operation 型が存在。前者は upsert_fact のみ対応、後者は全 op type を定義している。MVP の apply-patch が前者を使っているため schema が不完全 | `patch/model.go` を正として、extraction 側を統合 |
| **API 認証なし** | `/api/v1/*` は全て認証なしで公開。apply-patch は workspace_path を自由に受け取れる状態 | Phase 1 完了前に JWT/session 認証を導入 |
| **Fact schema の不整合** | docs/implementation/workspace/fact-schema.md では `action`, `decision`, `impact`, `context`, `evidence` フィールドが定義されているが、実際の YAMLFact struct はこれらを `description` に集約している。docs と実装の乖離 | Phase 2 で fact schema を拡張する際に統一。それまでは docs を「将来計画」として扱う |
| **go.sum の残留依存** | yaml.v3 のみが必要なはずだが go.sum に不要なエントリが残っている可能性 | `go mod tidy` で整理 |

### 中優先度（Phase 2 で対処）

| 負債 | 詳細 | 対処方針 |
|---|---|---|
| **packages/* 未分離** | workspace-core, patch-engine, schema, shared はすべて README のみ。実際の実装は apps/api 内に集まっており、パッケージ責務分離が形骸化している | Phase 2 以降で packages/* に実装を移行 |
| **Mock provider の keyword matching** | MockExtractionProvider は入力テキストのキーワードマッチングで fact を抽出する。テスト用途以外では使い物にならない | CodexExtractionProvider 実装で置き換え |
| **File System Access API 依存** | ブラウザ側のワークスペース attach は Chrome/Edge のみ対応。Firefox/Safari 非対応 | Electron ラッパーまたは CLI ヘルパー経由でのアクセスを検討 |
| **ストリーミング未実装** | AI 応答が全量返るまでユーザーが待つ必要がある | SSE または WebSocket 対応 |

### 低優先度（Phase 3 以降）

| 負債 | 詳細 |
|---|---|
| profiles/narratives 命名不統一 | docs 内で `profiles/` と `narratives/` が混在。正式名称未確定 |
| embeddings の位置づけ未確定 | CareerVault に embeddings/ を持つかどうか未決定 |
| sync 設計なし | マルチデバイス同期の設計が存在しない |

---

## 6. ロードマップ：3機能の実現計画

### 前提：Phase 1 MVP 完成（現在進行中）

3機能のロードマップは Phase 1 完了を前提とする。Phase 1 残り作業は次の通り。

```
Phase 1 残り（目安: 3-4週間）
  ├── AI 実統合（CodexExtractionProvider）
  ├── fact 整形 UI
  ├── profile 生成（API + UI）
  ├── export（API + UI）
  ├── 認証（JWT + ログイン画面）
  └── E2E 統合テスト
```

---

### Feature A: マルチユーザー対応

**目的**: 単一ユーザー前提の MVP から、複数ユーザーが独立してサービスを利用できる状態へ。

**設計原則の維持**: ユーザーアカウントはプラットフォーム DB（PostgreSQL）で管理するが、CareerVault（facts/profiles/exports）は引き続き各ユーザーのローカルに置く。DB はキャリア情報の正本にならない。

#### Phase A-1: 認証基盤（Phase 1 内で着手可）

| タスク | 詳細 |
|---|---|
| ユーザー認証 | JWT または session-based auth。メール/パスワード or OAuth（GitHub/Google）。Supabase Auth または自前 JWT |
| API 認証ミドルウェア | Go HTTP middleware で全 `/api/v1/*` を保護 |
| ログイン画面 | React: ログイン・登録フォーム |
| セッション管理 | user session と workspace attachment を別管理（auth-flow.md の方針に従う） |

```
CareerNESS の認証スコープ：
  ✅ AppServer セッション確立
  ✅ AI 利用権限の確認
  ❌ CareerVault ownership 移転（これはしない）
```

#### Phase A-2: プラットフォーム DB 導入

| タスク | 詳細 |
|---|---|
| DB 選定 | Supabase（PostgreSQL hosted）または Railway PostgreSQL。MVP では SQLite でも可 |
| スキーマ設計 | users, sessions, workspace_attachments, usage_logs テーブル（career data は入れない） |
| per-user workspace isolation | workspace_path を session に紐付け。他ユーザーの workspace_path を apply-patch が受け付けないようにする |
| rate limiting | ユーザーごとの AI 呼び出し制限（OpenAI cost 管理） |

#### Phase A-3: マルチデバイス対応（Optional）

| タスク | 詳細 |
|---|---|
| workspace 選択 UI | 複数ワークスペースの切り替え（本ローカルフォルダの選択）。ユーザーごとに最後に attach したワークスペースを記憶 |
| Git-backed sync（検討） | CareerVault を Git リポジトリとして管理し、GitHub/GitLab での同期を選択肢として提供。運営側が正本を保持しない形での sync |
| sync 方針文書化 | Local-first を壊さない sync の原則を ADR-004 として記録 |

#### 完了基準

- 複数ユーザーが同時にサービスを利用できる
- API が認証保護されている
- あるユーザーが別ユーザーの workspace にアクセスできない
- login = workspace ownership の誤結合が起きていない
- CareerVault はローカルのまま（DB に career data が入っていない）

---

### Feature B: AIキャリア分析

**目的**: 単純な fact 抽出から、ユーザーのキャリア全体を分析し、洞察を提供する機能へ。

**前提**: CodexExtractionProvider（real AI）が実装されていること。confirmed facts がある程度蓄積されていること。

#### Phase B-1: AI 実統合（Phase 1 内）

| タスク | 詳細 |
|---|---|
| CodexExtractionProvider 実装 | Codex App Server 経由で Claude API を呼び出し。ユーザー自身の OpenAI/ChatGPT アカウントを利用 |
| プロバイダー切り替え | ExtractionProvider interface を活用し、mock / Codex を環境変数で切り替え |
| ストリーミング対応 | AI 応答を SSE でフロントエンドへ stream |
| clarification question 表示 | AI が返した質問をチャット UI に表示し、ヒアリングを継続 |

#### Phase B-2: 多段階抽出・精緻化

| タスク | 詳細 |
|---|---|
| multi-stage extraction | 会話を重ねて facts を段階的に詳細化。初回: summary レベル → 2回目以降: action/decision/impact の深掘り |
| fact schema 拡張 | `action`, `decision`, `impact.metrics`, `context.constraints`, `evidence` フィールドを YAMLFact に追加（fact-schema.md との整合） |
| セマンティック重複検出 | 同一 project/period の fact が既存 facts に存在する場合に警告・統合提案 |
| confidence propagation | 複数の会話ターンをまたいで confidence を更新する仕組み |

#### Phase B-3: キャリア分析ダッシュボード

| 機能 | 詳細 |
|---|---|
| **キャリアタイムライン** | facts を period ベースで時系列表示。プロジェクト・役割の遷移を可視化 |
| **スキル分布分析** | tech_stack / tags の集計。時期ごとのスキル変化グラフ |
| **ロール別 fact マッピング** | 「バックエンドエンジニア」「テックリード」「プラットフォームエンジニア」など役割ごとに facts を分類 |
| **キャリアギャップ検出** | 期間が不明な facts、action のない facts、impact が未記入の facts を洗い出してユーザーに通知 |
| **強み・パターン抽出** | confirmed facts 全体から繰り返し登場する行動パターン・技術領域を AI が分析 |

#### Phase B-4: 目標ロール分析（JD 解析）

| 機能 | 詳細 |
|---|---|
| **JD インポート** | 求人票（テキスト）を貼り付けると、AI が要求スキル・経験を構造化 |
| **ギャップ分析** | JD の要求と自分の facts を照合。カバーできている・できていない領域をマッピング |
| **強調すべき facts の提案** | この JD に対して「どの facts を profile に含めるべきか」を AI が提案 |
| **profile 自動生成** | JD に基づいて facts を選別し、profile を自動生成。Patch proposal として提示 |

#### 完了基準

- real AI（Codex 経由）で facts が抽出できる
- 会話を重ねて facts が段階的に詳細化される
- キャリアタイムラインが表示される
- JD を貼り付けると、対応する facts のマッピングが表示される

---

### Feature C: AI面接支援

**目的**: 蓄積した facts を活用して、面接準備を構造的に支援する。

**前提**: Feature B の AIキャリア分析が一定程度完成していること（facts が十分に蓄積・確認済みであること）。

#### Phase C-1: 面接質問生成・回答フレームワーク

| 機能 | 詳細 |
|---|---|
| **STAR 形式回答生成** | facts の `action` / `decision` / `impact` / `context` を STAR (Situation-Task-Action-Result) フレームワークに変換。Patch proposal として export に保存 |
| **よく聞かれる質問への回答準備** | 「最も困難だった課題は？」「チームをどうリードしたか？」などの典型質問に対して、関連 facts を紐付けて回答案を生成 |
| **ポジション別質問予測** | 「シニアバックエンドエンジニア」「テックリード」など target role を指定すると、そのポジションで聞かれやすい質問を AI が予測 |
| **回答の強さ評価** | 生成した回答に対して、具体性・数値的根拠・影響の明確さの観点で評価し、改善提案 |

**実装方針**: 面接支援も patch proposal model に従う。AI が生成した回答案は必ず `exports/interview-prep/` への Patch として提示し、ユーザーが確認・承認してから保存。

#### Phase C-2: 模擬面接フロー

| 機能 | 詳細 |
|---|---|
| **模擬面接モード** | チャット UI に「面接モード」を追加。AI が面接官として質問を出す |
| **回答評価・フィードバック** | ユーザーの回答に対して AI がフィードバック（具体性、STAR の完成度、facts との整合性） |
| **面接後の fact capture** | 模擬面接中に出てきた新しいエピソードを AI が検出し、新規 fact の Patch proposal として提示 |
| **弱点エリアの特定** | 回答が弱かった質問カテゴリを記録し、強化すべき facts の収集を促す |

#### Phase C-3: 面接準備 export

| 機能 | 詳細 |
|---|---|
| **面接準備シート生成** | 企業・ポジション別の「面接チートシート」を export。使う facts・強調点・STAR 回答サマリーを含む |
| **一言 elevator pitch 生成** | 30秒・1分・3分のキャリアサマリーを facts から生成 |
| **技術面接準備** | tech_stack から「説明できるはずの技術」リストと、それに紐づく facts のマッピング |
| **想定 Q&A ドキュメント** | 想定質問と回答案のペア一覧を markdown で export。CareerVault の `exports/interview-prep/` に保存 |

#### 完了基準

- 蓄積した facts から STAR 形式の回答案が生成できる
- 面接モードで AI との模擬面接ができる
- 面接準備シートを markdown で export できる
- 面接で出てきた新エピソードが fact として capture できる

---

## 7. フェーズ統合ロードマップ

```
現在（2026-05）
│
├─ Phase 1 MVP 完成（~2026-06末）
│   ├── AI 実統合（CodexExtractionProvider）
│   ├── fact 整形 UI
│   ├── profile 生成
│   ├── export
│   └── 認証基盤（Feature A-1 と並行）
│
├─ Phase 2 マルチユーザー基盤（~2026-07末）  ← Feature A
│   ├── プラットフォーム DB（Supabase）
│   ├── per-user workspace isolation
│   ├── rate limiting
│   └── マルチデバイス sync 検討
│
├─ Phase 3 AI キャリア分析（~2026-09末）  ← Feature B
│   ├── multi-stage extraction
│   ├── fact schema 拡張（action/decision/impact）
│   ├── キャリアタイムライン UI
│   ├── スキル分布・ロール別マッピング
│   └── JD 解析・ギャップ分析
│
└─ Phase 4 AI 面接支援（~2026-12末）  ← Feature C
    ├── STAR 形式回答生成
    ├── ポジション別質問予測
    ├── 模擬面接モード
    ├── 面接後 fact capture
    └── 面接準備シート export
```

### 依存関係

```
Phase 1（MVP）
  └─→ Feature A（マルチユーザー）
        └─→ Feature B（AIキャリア分析）
              └─→ Feature C（AI面接支援）
```

Feature A と Feature B の一部（AI 実統合、multi-stage extraction）は並行着手可能。  
Feature C は Feature B で facts が十分に蓄積された状態を前提とするため、後続フェーズ。

---

## 8. アーキテクチャ変更ポイント（3機能実現時）

### マルチユーザー対応で追加が必要なもの

```
現在:
  Browser UI → Go API (no auth) → Local workspace

追加後:
  Browser UI
    ↓ (JWT / session token)
  Go API
    ↓ user session lookup
  Supabase (platform DB: users, sessions, usage)
    ↓ workspace_path resolved per user
  Local workspace (CareerVault, unchanged)
```

注意: Supabase には career data を入れない。platform metadata のみ。

### AIキャリア分析で追加が必要なもの

```
現在:
  /api/v1/conversations/message → MockExtractionProvider → Patch

追加後:
  /api/v1/conversations/message → CodexExtractionProvider
    → Codex AppServer (via user's OpenAI account)
    → Structured JSON extraction
    → multi-stage refinement loop
    → enriched Patch (action/decision/impact)
    → Patch proposal
```

分析 API として新規追加が必要:
- `GET /api/v1/analysis/timeline` — facts の時系列分析
- `GET /api/v1/analysis/skills` — スキル分布集計
- `POST /api/v1/analysis/jd` — JD テキスト解析・ギャップ分析

### AI面接支援で追加が必要なもの

```
新規 export type として追加:
  /api/v1/interview/questions    — ポジション別想定質問生成
  /api/v1/interview/star         — facts → STAR 回答生成
  /api/v1/interview/mock/start   — 模擬面接セッション開始
  /api/v1/interview/mock/message — 面接官 AI との対話
```

面接支援も patch proposal model に従い、生成した回答案・準備シートは  
必ず `exports/interview-prep/` への Patch として提示する。

---

## 9. 優先度判断の指針

### 今すぐやるべきこと（Phase 1 完成のため）

1. `go mod tidy` による go.sum 整理（30分、負債解消）
2. Patch struct の二重定義解消（`extraction/models.go` と `patch/model.go` の統一）
3. API 認証ミドルウェアの実装（セキュリティ上最優先）
4. CodexExtractionProvider の実装（mock 卒業）
5. E2E 統合テスト（extract → patch → apply の full loop）

### Feature A → B → C の順序を守る理由

- **マルチユーザーなしの AIキャリア分析** は単一ユーザーでも技術的には可能だが、認証がないまま公開すると API が無保護で AI コストが制御できない
- **AIキャリア分析なしの AI面接支援** は面接支援の回答品質が facts の量・質に依存するため、先に facts 収集・確認フローを磨く必要がある
- 各 Feature は前段の output を input として利用する構造

### Local-first 原則を壊さないチェックリスト

新機能追加時に毎回確認する観点:

- [ ] CareerVault の facts/profiles/exports を DB に書いていないか
- [ ] workspace 外のパスを AI が参照していないか
- [ ] ユーザー未承認の変更が自動 apply されていないか
- [ ] session 内の inferred state が hidden truth 化していないか
- [ ] export が fact の代替として扱われていないか

---

## 10. 参照ドキュメント

| ドキュメント | 内容 |
|---|---|
| `docs/product/vision.md` | プロダクトビジョン・Local-first の理由 |
| `docs/product/non-goals.md` | やらないことの明示 |
| `docs/architecture/system-overview.md` | 3層構成の全体像 |
| `docs/architecture/responsibility-boundary.md` | Browser/AppServer/Workspace/Cloud の責務 |
| `docs/architecture/platform-vs-vault.md` | DB と CareerVault の責務分離 |
| `docs/implementation/decisions/001-003` | ADR: workspace boundary, patch model, local-first |
| `docs/implementation/workspace/fact-schema.md` | Fact の詳細スキーマ（将来計画含む） |
| `docs/implementation/ai/extraction-specification.md` | extraction pipeline 仕様 |
| `docs/implementation/auth/auth-flow.md` | 認証フロー設計 |
| `docs/roadmap/roadmap.md` | オリジナル roadmap（Phase 0-4） |
| `IMPLEMENTATION_PROGRESS.md` | 実装進捗・技術的負債一覧 |
