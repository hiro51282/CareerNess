# Platform と CareerVault の責務境界

CareerNess は将来的に user account / subscription / auth などの platform 機能を持つ可能性がある。そのため DB を使う場面が生まれる。ただし重要なのは、「DB を使うこと」と「CareerVault の truth を DB に置くこと」は別の問題である点だ。

## 原則

**DB = platform metadata**  
**CareerVault = career truth**

この2つの責務分離を維持する。

## DB に置いてよいもの

将来的に DB (PostgreSQL / Supabase / RDS 等) に置いてよいもの：

- user account
- subscription / billing
- usage control / rate limit
- auth session
- audit metadata
- AppServer の transient session metadata

これらは「CareerNess というプラットフォームを動かすための情報」であり、ユーザーのキャリアの正本ではない。

## DB に置いてはいけないもの

以下は DB に置かない：

- canonical facts
- canonical profiles
- canonical exports
- workspace の hidden cache / mirror
- cloud-side truth として機能するキャリアデータ

CareerVault が canonical source である。DB はその代替にならない。

## なぜこの分離が必要か

CareerNess は「クラウド DB が truth を持つ SaaS」ではなく、「local workspace attached AI runtime」に近い設計を目指している。

便利さのために DB persistence を増やしていくと、次の drift が起きやすい：

- AppServer の session cache が hidden truth になる
- profile の再生成結果だけが DB に残り、workspace に戻らない
- export を DB に保存することで workspace との乖離が生まれる
- login state と workspace ownership が混同される

これらを防ぐため、DB の用途を platform concern に限定する。

## MVP のスタンス

DB 選定は MVP では未確定でよい。現時点では：

- workspace attach
- patch proposal
- patch review
- yaml apply

などの core runtime を優先する。DB が必要になるのは auth / subscription / rate limit を実装するフェーズから。

## 実装上の注意

- AppServer が career data を DB に書く API を作らない
- session に乗せた inferred facts を DB に永続化しない
- export 再生成結果を DB に保存して workspace の代わりにしない
- cloud sync convenience のために CareerVault の mirror を DB に持たない
