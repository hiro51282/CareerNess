# System Overview

CareerNess の基本構成はシンプルである。

```text
Browser UI
  ↓
Codex AppServer
  ↓
CareerVault (Local Workspace)
```

この構成の意図は、UI、AI orchestration、ユーザーデータの正本を分離することにある。特に重要なのは、CareerVault がクラウド上のデータストアではなく、ユーザーが所有する local workspace である点である。

## 全体像

### Browser UI

Browser は、ユーザーが AI と対話し、workspace を attach し、facts / profiles / exports を閲覧する入口である。UI は会話、確認、承認、編集のための層であり、キャリア情報の正本を長期保持する層ではない。

### Codex AppServer

AppServer は、認証済みユーザーセッションのもとで AI 呼び出しや workspace 接続を調停する。ここで重要なのは、AppServer 自体が CareerNess の正本データベースにならないことだ。AppServer の責務は orchestration であり、恒久的なキャリア保管ではない。

### CareerVault

CareerVault は、ユーザーのキャリア情報を保持する local workspace である。facts、profiles、exports、projects などのファイル群を人間が読める形式で格納し、AI が触れる範囲もこの workspace に限定する。

## この構成で守りたいこと

- キャリア情報の正本は workspace に置く
- Browser は見せる・入力を受けるが、正本を持たない
- AppServer は AI を使えるようにするが、ユーザーデータを所有しない
- Cloud は必要最小限の補助にとどめる

## データの流れ

典型的な流れは次の通り。

1. ユーザーが Browser でログインする
2. ユーザーが CareerVault を attach する
3. Browser が AppServer 経由で AI セッションを開始する
4. AI は attach 済み workspace の中だけを参照して質問や提案を行う
5. facts の追加や profile 更新は patch proposal として提示される
6. ユーザー承認後に workspace へ反映される
7. 必要に応じて export を生成する

## 設計上の判断

### なぜ Browser 直結ではなく AppServer を置くのか

AI 連携、認証、セッション管理、workspace attach の制御を UI から分離したいためである。Browser から直接あらゆる能力を扱わせると、責務が肥大化し、能力境界の説明もしにくくなる。

### なぜ Cloud DB 中心にしないのか

CareerNess は Local-first を重視しているため、キャリア情報の正本を cloud DB に置くと思想と矛盾しやすい。クラウド保存は便利だが、ownership と capability boundary が曖昧になりやすい。

### なぜ CareerVault をファイルベースで考えるのか

人間が直接読めること、Git やバックアップ運用と相性がよいこと、AI の編集対象を明確にしやすいことが理由である。構造化データであっても、完全に opaque な内部 DB に閉じる設計は優先しない。

## MVP 前提

MVP では次を優先する。

- 単一ユーザーの workspace attach
- facts / profiles / exports の責務分離
- AI の workspace-scoped 動作
- patch proposal ベースの更新

次は MVP では前提にしない。

- 組織全体の共有 workspace
- 複雑なクラウド同期
- 高度な multi-region / multi-tenant 設計
