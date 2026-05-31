# Auth Flow

CareerNess の auth は、ユーザーのキャリア正本を運営側で所有するための仕組みではない。ここで認証するのは「誰が AppServer と AI 機能を使うか」であり、「誰が workspace の owner か」を補助するためである。

## Goal

- Browser と AppServer のセッションを成立させる
- ユーザー自身の OpenAI / ChatGPT アカウント利用を前提にする
- workspace attach 操作をユーザー起点で明示化する

## Scope

auth flow が扱うもの:

- CareerNess UI へのログイン
- AppServer セッション開始
- AI 利用権限の確認
- workspace attach の明示的開始

auth flow が扱わないもの:

- CareerVault 正本のクラウド保管
- キャリア事実の所有権移転
- multi-user collaboration

## Core Flow

```text
Browser access
  ↓
User login
  ↓
AppServer session established
  ↓
User selects local workspace
  ↓
Workspace attach granted
  ↓
AI-assisted session starts
```

## Principles

- ログイン済みでも workspace は自動 attach しない
- attach しても workspace の owner はユーザーのまま
- AppServer セッション終了後、workspace ownership は何も移らない
- OpenAI 利用条件はユーザー側アカウント状態に依存する
- login は AI access / AppServer session / usage control のためであり、workspace ownership を作らない

## Session Binding

認証だけでは不十分で、次の 2 つを別で持つ。

- user session
- workspace attachment

この分離により、同一ログイン中でも workspace を切り替えたり、attach を解除したりできる。
同時に、`login = workspace owner` という誤読も避けられる。

## Minimal Retention

AppServer 側で保持する情報は最小限にとどめる。

- 認証状態
- セッション識別子
- AI 実行に必要な最小メタデータ

保持しないもの:

- CareerVault の永続コピー
- facts / profiles / exports の正本

## Failure Cases

- ログイン成功だが AI アカウント利用不可
- セッション有効だが workspace 未 attach
- attach 済みだが patch apply 権限未承認

これらはそれぞれ別エラーとして扱う。認証と workspace access を一つの成功/失敗に潰さない。

## Prohibited Patterns

- ログイン直後に過去 workspace を自動再接続する
- サーバー側で workspace 内容をキャッシュ正本化する
- auth 成功を理由に approval を省略する
- login 情報だけで local apply を許可する

## workspace_attachments テーブル設計指針（2026-05-31 確定）

Reviewer レビュー（e58571e）を経て以下を確定した。

`workspace_attachments` テーブルに持たせてよいカラムと禁止カラムを明示する。

```sql
-- OK: このカラムのみ持たせる
workspace_attachments (
  id              UUID PRIMARY KEY,
  user_id         UUID NOT NULL REFERENCES users(id),
  session_id      UUID NOT NULL,
  workspace_root_hash  TEXT,  -- ディレクトリ識別用の不透明なハッシュのみ
  attached_at     TIMESTAMPTZ NOT NULL,
  detached_at     TIMESTAMPTZ
)

-- NG: 以下は持たせない
--   last_seen_files[]     ← NG: ファイルリストはキャリアデータの proxy になる
--   cached_fact_ids[]     ← NG: fact の参照は career data の影
--   workspace_path        ← 程度による: フルパスよりハッシュを優先
--   last_synced_at        ← NG: sync を前提とした設計に滑り込む
```

この設計により、DB は「どのユーザーがどのセッションで workspace を attach していたか」を記録するに留まり、workspace の内容・構造・キャリアデータには触れない。

## per-user workspace isolation の実装要件（P0）

認証が完成しても isolation が別途実装されなければ、あるユーザーが別ユーザーの `workspace_path` を `apply-patch` に渡すことができてしまう。

以下を Phase 1 認証タスク着手前に設計・実装する（認証とは独立した課題）：

1. workspace_path を session に紐付けるタイミング（attach 操作時）
2. `apply-patch` が受け付けるパスを session の `workspace_root` 以下に限定する validation
3. ユーザー A の session token でユーザー B の `workspace_path` が `apply-patch` に通らないことの保証

## session 内 inferred state の扱い（確定）

会話が長くなった場合の inferred state（まだ patch になっていない AI の中間出力）はメモリのみに置く。

- session 切断で消える → **これが正しい挙動**
- プロダクトはユーザーのキャリア情報に責任を持たない（プロダクトオーナー追加前提）
- 消えたくない情報はユーザーが patch proposal を承認して workspace に保存する

## Open Questions

- ローカルクライアント連携の認可方式は未確定（Phase 1 完了前に決定が必要）
- offline-first 時の session expiry 振る舞いは未確定（Phase 1 完了前に決定が必要）
- ユーザーの OpenAI API key を AppServer がどのように受け取るかの設計が未確定（ADR-004 の残課題）
