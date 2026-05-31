# Planner Response — Reviewer Report への再評価

作成者: PM エージェント (Claude Sonnet 4.6)  
作成日: 2026-05-31  
参照: `docs/reviewer-report.md`（ブランチ agent-review / コミット e58571e）

## プロダクトオーナーから確認された追加前提

以下の前提が本再評価の判断基準として追加された。

1. プロダクトはユーザーのキャリア情報に責任を持たない
2. Local-first はその原則を実現するための設計である
3. LLM 利用のための API キーおよび API 利用料金はプロダクト側では提供しない
4. ユーザー自身が利用する AI サービスの利用料金を負担する
5. **AppServer はユーザー自身の OpenAI アカウント（Codex 経由）を利用する構造である**
6. CareerNESS 運営は AI 利用料の請求主体にならない

追加前提 5 が、Reviewer の [CRITICAL] 指摘に直接回答する情報であるため、以下の再評価はこれを軸に行う。

---

## 1. Reviewer の指摘で受け入れるもの

### 1-1. 認証戦略の実装空白（受け入れ）

auth-flow.md の原則設計（login ≠ workspace ownership）は正しいという評価に同意する。加えて Reviewer が指摘した以下は全て有効な課題である。

- JWT クレーム設計が未定義
- workspace_attachments テーブルに `last_seen_files[]` や `cached_fact_ids[]` を持たせると Local-first 原則が崩れる
- token refresh ポリシーが未定義

特に workspace_attachments テーブルの NG 例 / OK 例の提示は具体的かつ正確。OK 例（`workspace_root_hash`、`attached_at`、`detached_at` のみ）を実装方針として採用する。

### 1-2. per-user workspace isolation の実装詳細が不足（受け入れ）

`POST /api/v1/apply-patch` が現在 `workspace_path` を任意のパスで受け取っている事実、および「認証 ≠ isolation」の整理は正確。session の workspace_root 外への apply-patch を拒否する validation が認証とは独立して設計・実装される必要があることを計画に明示する。

### 1-3. fact schema の minimal fix（受け入れ）

Phase 2 での fact schema 変更が confirmed facts のマイグレーションを引き起こすリスクは有効。

Reviewer の推奨案「Phase 1 の時点で `action / decision / impact` の3フィールドを空で追加する」を採用する。空フィールドの場合は `description` にフォールバック表示する設計で互換性を保てる。これはコード変更を伴うため、Phase 1 残りタスクへ追加する。

### 1-4. packages/* 分離スプリント（受け入れ）

Phase 2 着手前に packages 分離スプリントを入れる推奨を受け入れる。特に「Feature A で認証ミドルウェアと workspace isolation を追加する際に混乱構造では影響範囲が不明確」という理由は説得力がある。

最小スコープ（schema → patch-engine → workspace-core の順に移行、apps/api は orchestration + handler のみ残す）をロードマップに明示する。

### 1-5. Phased-approval（session buffer 型）UX（受け入れ）

面接支援の「模擬面接中に都度 patch を承認するフロー」の UX 問題に対し、session buffer 型（面接後にまとめて承認）を採用する。ADR-002（patch-oriented editing）の精神を維持しつつ会話テンポを損なわない設計として妥当。

ただし実装上の注意として Reviewer が示した「session buffer は AppServer のメモリ内のみ（DB に書かない）」「セッション終了 = buffer 消去」は追加前提 1（プロダクトはユーザーのキャリア情報に責任を持たない）とも整合する。

### 1-6. Local-first チェックリストの PR テンプレート組み込み（受け入れ）

箇条書きのまま docs に置くより、`.github/pull_request_template.md` に埋め込む方が機能する。CI 自動検証できるもの（workspace 外パスへの apply-patch が 400 を返すか）はテスト化する提案も採用する。

### 1-7. session 内 inferred state の扱い（受け入れ）

「まだ patch になっていない AI の中間出力をどこで持つか未規定」は有効な指摘。会話が長くなった場合の挙動（session 切断で消える）を「正しい挙動」として明示的にドキュメント化する必要がある。これは追加前提 1（プロダクトはキャリア情報に責任を持たない）からも「消えてよい」が正解。

---

## 2. Reviewer の指摘で誤解しているもの

### 2-1. [CRITICAL] Codex AppServer の矛盾評価（部分的に誤解あり）

**Reviewer の指摘**: 3つのドキュメントが3方向に発散しており、候補 A（Claude API 直接）/ B（OpenAI API 直接）/ C（独自中間サーバー）のどれかを選んで確定せよ。

**再評価**:

追加前提 5 により、この問いに対する答えが確定した。

> AppServer はユーザー自身の **OpenAI アカウント（Codex 経由）** を利用する構造である

これにより以下が確定する。

| ドキュメント | 評価 | 意味 |
|---|---|---|
| CLAUDE.md「Codex AppServer 経由」 | ✅ 正しい | `apps/api`（Go API）がユーザーの OpenAI Codex を呼ぶ中継層 |
| auth-flow.md「ユーザー自身の OpenAI / ChatGPT アカウント利用」 | ✅ 正しい | 追加前提 3〜5 と完全整合 |
| extraction-specification.md の `github.com/anthropic-ai/sdk-go` | ❌ 誤り | **設計ドキュメントの誤記**。Anthropic SDK を使う前提は追加前提と矛盾する |

つまり「矛盾している」という Reviewer の診断は正しいが、解決策は候補 B（OpenAI API を直接呼ぶ）が正解。候補 A（Claude API 直接）は追加前提により採用しない。候補 C（独自中間サーバーが実在する）は不要。

**Reviewer が誤解している点**: これを「P0 の判断事項として残る未解決問題」として扱っているが、追加前提により既に解決されている。P0 から降格する。

### 2-2. rate limiting の目的認識（誤解あり）

**Reviewer の指摘**: 「認証なしで AI を公開すると API コストが制御できない」。section 6 Feature A-2 の rate limiting をコスト管理のための機能として位置付けている。

**再評価**:

追加前提 3・4・6 により、**CareerNESS 運営は AI 利用コストの管理責任を負わない**。ユーザーが自分の OpenAI API key を使い、自分のコストを負担する。

rate limiting が必要な理由は「OpenAI コスト管理」ではなく、**AppServer 自体の abuse 防止**（認証を突破した不正利用者が AppServer を踏み台に使うことの防止）に限定される。

この整理により:
- per-user の AI 呼び出し回数制限（OpenAI コスト保護のための制限）は CareerNESS の責務外
- AppServer エンドポイントへの rate limiting は依然必要（abuse 防止として）

Feature A-2 の「rate limiting」タスクの目的説明を「OpenAI cost 管理」から「AppServer abuse 防止」に修正する。

---

## 3. 追加前提を踏まえた修正が必要な箇所

### 修正対象: `docs/agent-plan.md`

#### 修正 1: Section 1（プロジェクト概要）へ前提追記

以下の前提をプロジェクト概要テーブルに追加する。

```
| キャリア情報の責任 | プロダクトはユーザーのキャリア情報に責任を持たない。正本の管理はユーザー責任 |
| AI 利用料金 | ユーザー自身の OpenAI アカウントを使用。CareerNESS 運営は請求主体にならない |
```

#### 修正 2: Section 6 Feature A-2（マルチユーザー基盤）の rate limiting 目的

変更前: 「ユーザーごとの AI 呼び出し制限（OpenAI cost 管理）」  
変更後: 「AppServer abuse 防止のための rate limiting（AI コスト管理は CareerNESS の責務外）」

#### 修正 3: Section 6 Feature B-1（AI 実統合）の実装方針

変更前: 「Codex App Server 経由で Claude API を呼び出し。ユーザー自身の OpenAI/ChatGPT アカウントを利用」（Claude API と OpenAI が混在）  
変更後: 「ユーザー自身の OpenAI アカウント（Codex 経由）を使って AI 呼び出しを実行。`apps/api` が中継層として機能。`extraction-specification.md` の Anthropic SDK 実装例は誤記であり、OpenAI SDK を使う実装に差し替える」

#### 修正 4: Section 5 技術的負債に追記

`extraction-specification.md` における `github.com/anthropic-ai/sdk-go` 使用は設計ドキュメントの誤記として負債リストに追加する。

```
| extraction-specification.md の Anthropic SDK 誤記 | 追加前提と矛盾する。OpenAI SDK を使う実装例に差し替えが必要 | 高優先度（Phase 1 AI 実統合前） |
```

#### 修正 5: Section 7（フェーズ統合ロードマップ）に packages 分離スプリントを追加

```
Phase 1 完成後・Phase 2 着手前:
  packages 分離スプリント（~1週間）
    ├── packages/schema → YAMLFact, Patch 型の一元化（二重定義解消）
    ├── packages/patch-engine → validator, applier 移行
    └── packages/workspace-core → file I/O, attach 管理 移行
```

#### 修正 6: Phase 1 残りタスクに fact schema minimal fix を追加

```
Phase 1 残り（更新版）
  ├── AI 実統合（OpenAI Codex 経由の CodexExtractionProvider）
  ├── fact schema minimal fix（action/decision/impact を空フィールドとして追加）  ← 追加
  ├── fact 整形 UI
  ├── profile 生成（API + UI）
  ├── export（API + UI）
  ├── 認証（JWT + workspace isolation + ログイン画面）
  └── E2E 統合テスト
```

### 修正対象: `docs/reviewer-handoff.md`（次回 Reviewer 引継ぎがある場合）

U-1（Codex AppServer 仕様の乖離）は追加前提で解決済みのため、ステータスを「解決済み」に更新する。

### 修正不要: 設計ドキュメント（コード変更禁止のため記録のみ）

- `docs/implementation/ai/extraction-specification.md` の Anthropic SDK 実装例 → 次の実装フェーズで OpenAI SDK 実装例に差し替える必要あり（コード変更禁止のため今回は修正しない）

---

## 4. ロードマップ変更の要否

### フェーズ構造変更: **不要**

A → B → C の順序と理由付けは追加前提後も妥当である。

- Feature A が先（認証なしで AppServer を公開すると abuse される）
- Feature B が先（facts の量と質が面接支援の品質を決める）

### タスク内容・記述の修正: **必要**

| 箇所 | 変更内容 |
|---|---|
| Phase 1 残りタスク | fact schema minimal fix（`action/decision/impact` 追加）を追加 |
| Phase 1 AI 実統合 | 「OpenAI Codex 経由」として方向性を確定 |
| Phase A-2 rate limiting | 目的を「コスト管理」→「abuse 防止」に変更 |
| Phase 1〜2 の間 | packages 分離スプリントを明示的に追加 |
| Feature C UX | Phased-approval（session buffer 型）を正式採用として記載 |

### 優先度変更: **あり**

| 事項 | 変更前 | 変更後 | 理由 |
|---|---|---|---|
| Codex AppServer 定義 | P0（未解決） | **解決済み**（P0 から除去） | 追加前提 5 で確定 |
| workspace isolation 実装設計 | P1 | **P0 に格上げ** | 認証と独立して設計が必要。認証完成後に isolation が漏れやすい |
| fact schema minimal fix | P1 | **P1 維持（Phase 1 内に追加）** | マイグレーションコスト回避のため早期対処 |
| extraction-specification.md 誤記修正 | 未記載 | **P0 に追加** | Phase 1 AI 実統合の実装方針を誤らせる誤記 |

### 更新後の P0 判断事項

| # | 事項 | 対処期限 |
|---|---|---|
| P0-1 | workspace isolation の実装設計（apply-patch の session-scope validation） | Phase 1 認証タスク着手前 |
| P0-2 | extraction-specification.md の Anthropic SDK 誤記を OpenAI SDK に差し替え | Phase 1 AI 実統合着手前 |
| P0-3 | workspace_attachments テーブルの OK 例をドキュメントとして確定（Reviewer OK 例を採用） | Phase A-2 着手前 |

---

## 再評価サマリー

| Reviewer 指摘 | 再評価結果 |
|---|---|
| [CRITICAL] Codex AppServer 矛盾 | 追加前提 5 で解決。正解は OpenAI Codex 経由。extraction-specification.md が誤記 |
| 認証戦略の実装空白 | 受け入れ。JWT 設計・workspace_attachments 設計・token refresh を Phase 1 前に確定する |
| workspace_attachments テーブル設計リスク | 受け入れ。OK 例（workspace_root_hash のみ）を方針として採用 |
| fact schema minimal fix | 受け入れ。Phase 1 残りタスクに追加 |
| packages 分離スプリント | 受け入れ。Phase 1 完成後・Phase 2 前に挿入 |
| Phased-approval UX | 受け入れ。Feature C の設計方針として採用 |
| rate limiting の目的 | 部分的に修正。「コスト管理」ではなく「abuse 防止」が正しい目的 |
| session 内 inferred state | 受け入れ。「session 切断で消えてよい」を明示ドキュメント化する |
| Local-first チェックリスト → PR テンプレート | 受け入れ |
| Git-backed sync は ADR-004 で先に判断 | 受け入れ。Phase A-3 着手前に ADR-004 として記録する |
