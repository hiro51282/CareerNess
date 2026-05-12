# Responsibility Boundary

CareerNess では、責務と能力境界を文書として固定することを重視する。特に重要なのは、Browser、AppServer、Workspace、Cloud のどこが何を担い、何を担わないかを明示することである。

## Browser

### やること

- ユーザーとの会話 UI を提供する
- workspace attach の入口を提供する
- facts / profiles / exports の閲覧と確認を行う
- patch proposal の確認と承認操作を受ける
- AI の提案内容を人間が読める形で提示する

### やらないこと

- キャリア情報の canonical source になること
- attach されていないローカル領域を独自に探索すること
- 承認なしで facts を確定保存すること
- 長期保管前提の巨大ストレージになること

### 保持するもの

- 一時的な UI state
- セッションに必要な最小限の情報
- 現在開いている conversation context

### 保持しないもの

- ユーザーキャリアの正本
- 長期保存前提の facts / profiles / exports
- 端末全体に対する権限

## AppServer

### やること

- 認証とセッション管理
- AI orchestration
- Browser と Workspace の仲介
- workspace-scoped capability の適用
- patch proposal の生成補助

### やらないこと

- ユーザーキャリアの正本保管
- unrestricted file access の提供
- unrestricted shell access の提供
- ユーザー承認を飛ばした自動確定

### 保持するもの

- セッション情報
- 最小限の運用ログ
- 一時的な AI 実行コンテキスト

### 保持しないもの

- CareerVault 全体の永続コピー
- profile/export の恒久保存
- attach されていない workspace の内容

## Workspace

Workspace は CareerVault を指す。CareerNess におけるキャリア情報の正本はここにある。

### やること

- facts を構造化データとして保持する
- profiles を見せ方ごとの派生表現として保持する
- exports を提出・共有用出力として保持する
- projects / tags / role などの補助情報を保持する
- AI が触れてよい範囲を明示する

### やらないこと

- OS 全体のファイル置き場になること
- AI が無制限に横断探索できる領域になること
- クラウド依存の必須コンポーネントになること

### 保持するもの

- ユーザーが所有する career facts
- 人間が読める profile と export
- 設計上必要なメタデータ

### 保持しないもの

- 端末全体の unrelated data
- 運営側が所有する hidden canonical state

## Cloud

Cloud は Browser と AppServer を支える補助的な層であり、CareerVault の代替ではない。

### やること

- 認証基盤の提供
- AI 利用に必要な接続補助
- セッション成立に必要な最小限のサーバー機能

### やらないこと

- CareerVault の正本化
- 巨大なユーザーデータ保管庫化
- 全キャリア情報の継続収集
- enterprise data lake 化

### 保持するもの

- 最小限の認証関連情報
- 運用上必要な最小限のメタデータ

### 保持しないもの

- ユーザーの全 facts / profiles / exports の恒久保存
- attach された workspace の常時ミラー

## なぜこの分離が必要か

責務分離が曖昧になると、次の問題が起きる。

- Browser が便利さのために正本を持ち始める
- AppServer が気づかないうちにデータ保管責務を持つ
- Workspace boundary が形骸化する
- AI が「どこまで見てよいか」が不明確になる

CareerNess では、便利さより先に ownership と boundary を固定する。その上で不足する UX を改善する。

## 重要な禁止事項

- attach していないパスを AI の入力に混ぜない
- facts と profile を同一ファイルで混在させない
- export を正本扱いしない
- Cloud 側に無自覚なデータ retention を追加しない
