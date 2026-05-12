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

## Open Questions

- ローカルクライアント連携の認可方式は未確定
- offline-first 時の session expiry 振る舞いは未確定
