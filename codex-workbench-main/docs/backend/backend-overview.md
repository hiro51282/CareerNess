# Backend Overview

CareerNess の backend は、巨大な業務基盤ではなく、Browser と Workspace の間で AI を安全に使うための薄い orchestration 層として設計する。

## Backend の主な役割

- AI orchestration
- auth / session 管理
- workspace attach の調停
- patch proposal 生成補助
- cloud-side minimalism の維持

## AI Orchestration

backend は、Browser からの要求を受けて AI セッションを成立させる。重要なのは、backend 自体がキャリア正本を所有しないことだ。

backend が扱うのは次のような責務である。

- 現在 attach されている workspace context の取り扱い
- AI へ渡す instruction boundary の付与
- fact extraction / profile generation / export drafting の実行補助

## Auth / Session

CareerNess はユーザー自身の OpenAI / ChatGPT アカウント利用を前提とするため、backend は account / session の仲介を担う。

ただし、ここでの session は「AI 利用のための接続状態」であって、「キャリアデータ保管のためのサーバーセッションストア拡張」を意味しない。

## Cloud-side Minimalism

backend はクラウド側機能を持つが、可能な限り薄く保つ。

意味すること:

- CareerVault の恒久コピーを持たない
- export をサーバー保管前提にしない
- ユーザーデータ retention を最小限にする

## User Data Non-retention

CareerNess では、運営側がユーザーの career data を蓄積するほど価値が高まるモデルを目指さない。そのため backend も、データ集約装置ではなく、接続と制御のための層として考える。

これは次を意味する。

- facts / profiles / exports の canonical storage は backend に置かない
- ログ収集は運用最小限にとどめる
- 将来的な分析用途のために暗黙にデータ保存しない

## Backend がやらないこと

- unrestricted local access の代行
- CareerVault のクラウド同期エンジン化
- enterprise 向け巨大ワークフロー基盤化
- ユーザーデータを使った内部学習前提の設計

## MVP の backend 前提

- 単純なセッション管理
- 明示的な workspace attach
- facts / profiles / exports に対する操作の経路制御
- 最小限の observability

## 未確定事項

- AppServer の具体スタックは未確定
- session log の粒度は未確定
- 将来的な background job の必要性は未確定
