# CareerNess

Local-first AI支援キャリア構造化プラットフォーム。

---

# Overview

CareerNess は、AIとの対話を通じてキャリア情報を構造化・再解釈・管理するための Local-first プラットフォームです。

従来の転職活動では、媒体や目的ごとに経歴情報を何度も編集し直す必要があります。

例えば：

* 職務経歴書
* スカウト返信
* 面接用要約
* エージェント向けフォーマット
* role別レジュメ
* 技術特化プロフィール

CareerNess は、経歴情報を「文章」ではなく「構造化データ」として扱うことを目指します。

目的は AI に全部を書かせることではなく、ユーザー自身のキャリアを：

* 掘り起こす
* 構造化する
* 複数のプロフィール表現（見せ方）へ変換する
* キャリア情報の正本として維持する
* 文脈に応じて再利用する

ための支援を行うことです。

---

# なぜ CareerNess を作るのか

転職活動には、構造的な問題があります。

同じキャリアであっても：

* 応募媒体
* 企業
* ポジション
* 技術スタック
* マネジメント比率

によって「求められる見せ方」が変わります。

その結果、実際に行っていた役割と、市場から見える役割がズレることがあります。

例えば：

* 技術選定
* CI/CD改善
* 開発プロセス整備
* cross-team coordination
* platform改善
* 開発標準化

のような役割を担っていたとしても、レジュメ上では特定の実装業務だけが強調され、本来の役割や価値が十分に伝わらない場合があります。

CareerNess は、そのズレを埋めるための試みです。

---

# Philosophy

## Local-first

キャリア情報の ownership はユーザー側にあるべきだと考えています。

CareerNess は可能な限りローカル保存を前提とし、運営側でユーザーデータを保持しない方向を目指します。

あなたのキャリア情報はクラウドへ保存されず、ローカル workspace 内へ保存されます。

一方で、これは「ローカルファイルを失うとデータも失われる」ことを意味します。

CareerNess はユーザーによる ownership を重視するため、バックアップやファイル管理もユーザー責任となる想定です。

---

## キャリア情報と見せ方を分離する

CareerNess は以下を分離して扱います。

* キャリアの事実情報
* プロフィール表現（見せ方）
* export formats

職務経歴書やプロフィールは「正本」ではなく、構造化されたキャリア情報から生成される派生データとして扱います。

---

## OpenAI Account Usage

CareerNess は Codex AppServer を利用して AI 機能を提供する想定です。

そのため、AI利用にはユーザー自身の OpenAI / ChatGPT アカウントが利用されます。

CareerNess 運営側が API Key を保持・配布する構成は想定していません。

つまり：

* AI利用制限
* 利用可能モデル
* 利用量制限

は、ユーザー自身の OpenAI アカウント状態に依存します。

軽い利用であれば ChatGPT Go プラン程度でも十分動作可能な構成を目指しています。

---

## AI-assisted, not AI-owned

AI は：

* 質問する
* 掘り起こす
* 構造化する
* プロフィール表現を提案する

ための存在です。

ユーザーのキャリアを「所有」するものではありません。

---

## Workspace-scoped AI

AI がアクセスできる範囲は、attach された workspace に限定されます。

CareerNess は unrestricted local AI を目指しません。

---

# Features

現在想定している機能：

* AIによるキャリアヒアリング
* career facts 抽出
* プロフィール生成
* レジュメ export
* Local workspace 管理
* role別プロフィール管理
* Local-first storage
* AI workspace boundary 制御

---

# Architecture Overview

```text
Browser UI
↓
Codex AppServer
↓
CareerVault (Local Workspace)
```

---

# Security Model

CareerNess は、明示的な workspace boundary を前提として設計されます。

現在の設計思想：

* Workspace-scoped access
* unrestricted shell access を持たない
* User-owned data
* Local-first storage
* 最小限の cloud-side retention

---

# Workspace Structure

例：

```text
CareerVault/
 ├── facts/
 ├── narratives/
 ├── exports/
 ├── projects/
 └── embeddings/
```

---

# Design Priorities

CareerNess では、特に以下の設計を重要視しています。

## Workspace Boundary

AI が「どこまで読めるか」「どこまで書けるか」を明示的に制御します。

CareerNess は unrestricted local AI を前提とせず、attach された workspace を中心に動作する設計を目指します。

---

## Structured Career Model

CareerNess は「文章」を直接編集するのではなく、構造化されたキャリア情報を中心に扱います。

重要なのは：

* キャリアの事実情報
* プロフィール表現
* export形式

を分離して管理することです。

---

## AI Capability Design

AI は自由にシステムを操作する存在ではなく、明示的に許可された範囲の中で動作します。

現在想定している設計：

* workspace-scoped access
* patch proposal based editing
* direct overwrite の抑制
* unrestricted shell access を持たない
* ユーザー承認前提の変更フロー

---

## Responsibility Boundaries

CareerNess は責務分離を重視します。

例：

* Browser UI
* Codex AppServer
* Local Workspace
* Cloud Service

は、それぞれ異なる責務を持ちます。

---

## Human-readable Design

CareerNess は自然言語ベースの設計書を重視しています。

これは：

* 人間による理解
* AI coding agent との協調
* 設計意図の維持
* 実装暴走の抑制

を目的としています。

---

# Documentation

## Product

* Vision
* Non-goals

## Architecture

* System overview
* Responsibility boundaries

## Workspace

* Workspace layout
* Data model

## AI

* AI behavior
* Prompting strategy

## Security

* Security model

## Deploy

* Deployment strategy

## Roadmap

* MVP roadmap
* Future plans

---

# Non-goals

CareerNess は以下を目的にしていません。

* GPU inference infrastructure hosting
* Local LLM 必須化
* enterprise HR platform 化
* 全キャリアデータの cloud保存
* unrestricted autonomous AI agents

---

# Current Status

現在は architecture / prototype phase です。

主に以下を進めています：

* workspace design
* structured data modeling
* AI boundary design
* local-first architecture
* AI-assisted workflow experimentation

---

# License

License TBD.
