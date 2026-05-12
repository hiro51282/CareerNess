# Security Model

CareerNess の security は、暗号や認証だけでなく、AI capability boundary をどう設計するかに強く依存する。この文書では、CareerNess が守るべき最小限の security model を定義する。

## Core Principles

- capability boundary を明示する
- Local-first を前提にする
- user-owned data を守る
- workspace attachment を明示的に行う
- no unrestricted file access を貫く

## Capability Boundary

CareerNess における中心的な security 課題は、「AI がどこまで見て、どこまで変えられるか」である。

そのため、能力境界を次のように扱う。

- AI が触れる対象は attach された workspace に限る
- AI の更新は patch proposal を基本にする
- facts への確定反映には user approval を要求する
- AppServer は unrestricted shell/file capability を前提にしない
- MVP では attach workspace 外 capability を追加しない

Security 上、これは単なる UX 制約ではない。能力境界が曖昧だと、ユーザーは AI に何を渡したのかを理解できなくなる。

MVP で明示的に scope 外とするもの:

- clipboard
- browser history
- screenshot capture
- arbitrary local files
- global desktop search
- background collection of unrelated files

## Local-first

Local-first は、privacy だけでなく blast radius を小さくする設計でもある。正本が local workspace にあることで、クラウド側の漏えい・誤保持・過剰 retention の影響を下げられる。

ただし Local-first は万能ではない。

- ローカル消失時のバックアップ責任はユーザー側に寄る
- 端末自体が安全であることは別問題である

CareerNess はこのトレードオフを隠さない。

## User-owned Data

CareerNess は、ユーザーのキャリア情報を運営側が所有・常時保持する前提を取らない。これは法務的な ownership の話だけではなく、実際にどこへ保存されるかの話でもある。

意味すること:

- facts / profiles / exports の正本は workspace にある
- cloud-side retention は最小限に抑える
- ユーザーは自分のファイルを直接確認・バックアップできる
- login や cloud session は workspace ownership を移さない

## Workspace Attachment

workspace attachment は security model の起点である。どの workspace が AI の対象かが曖昧なままでは、access control を説明できない。

要求事項:

- ユーザーが対象 workspace を明示的に attach する
- attach 状態が UI でわかる
- attach されていない領域は対象外とする

## No Unrestricted File Access

CareerNess では、AI にローカルファイル全域への unrestricted access を与えない。

禁止の理由:

- キャリアと無関係な私的データまで探索されうる
- workspace boundary が意味を失う
- 人間の理解可能性が下がる

したがって、「ローカルで動くから安全」とは考えない。重要なのは local であることではなく、scope が明示されていることである。

## Threats CareerNess が避けたいもの

- AI が unrelated local files を読むこと
- AI が hallucinated fact を正本へ混入すること
- Cloud に意図せず career data が残り続けること
- export が唯一の正本になって修正不能になること
- session memory が hidden database 化すること

## Security Responsibilities

### Browser

- attach 状態を明示する
- patch review を可能にする
- ユーザーに scope を誤認させない

### AppServer

- capability boundary を壊さない
- 最小限のセッション情報だけを扱う
- ユーザーデータ retention を拡大しない
- patch proposal を truth そのものとして保存し続けない

### Workspace

- 正本と派生物を分離する
- 人間が直接監査できる形式で持つ
- proposed / pending / confirmed の状態を session 外でも表現できる

## 非目標

- ゼロトラスト企業基盤の完全実装
- 組織監査向けの巨大ログシステム
- 全脅威をクラウド側で解決する設計

## 未確定事項

- session log の保存期間は未確定
- 将来的な encrypted sync を公式に扱うかは未確定
