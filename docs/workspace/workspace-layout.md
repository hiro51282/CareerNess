# Workspace Layout

CareerVault は、CareerNess におけるユーザー所有の local workspace である。ここでは、AI が触れる対象と、人間が読み書きするキャリア資産の置き場を明確にする。

MVP では、複雑な内部 DB よりも directory と markdown / structured file を中心にした layout を優先する。

## 例

```text
CareerVault/
├── facts/
├── profiles/
├── exports/
├── projects/
├── taxonomy/
├── inbox/
├── attachments/
└── meta/
```

README では `narratives/` の例もある。MVP では `profiles/` を推奨しつつ、既存表現との整合は未確定として扱う。

## Directory Responsibilities

### `facts/`

キャリアの事実情報を置く。CareerNess の正本に最も近い層であり、AI 生成文ではなく、確認可能な事実を構造化して保存する。

入れるもの:

- 職歴 facts
- 成果 facts
- 技術スタック facts
- responsibility / decision / impact facts
- interview や対話から抽出した確認済み事項

入れないもの:

- 用途別の長文プロフィール
- 応募先向けに最適化した言い回し
- 最終提出用ドキュメント

### `profiles/`

facts を元にした見せ方の層。たとえば backend 重視、platform 重視、engineering manager 寄りなど、文脈ごとの narrative を置く。

入れるもの:

- role 別プロフィール
- 媒体別プロフィール
- 応募方針に応じた要約

入れないもの:

- 事実の正本
- 応募提出そのものの最終形式

### `exports/`

外部提出・共有のための出力を置く。PDF、markdown resume、プレーンテキスト、媒体貼り付け用フォーマットなどが対象である。

入れるもの:

- 提出用レジュメ
- スカウト返信用整形テキスト
- 共有用サマリ

入れないもの:

- 正本 facts
- profile の設計意図

### `projects/`

プロジェクト単位の補助データを置く。facts と重なる部分はあるが、ここでは案件・所属・期間・チーム構成など、プロジェクト文脈を見やすく管理する。

### `taxonomy/`

tags、roles、skill categories など、再利用する分類語彙を置く。CareerNess では自由記述を完全に禁止しないが、分類語彙が散らばると profile 生成時に揺れが出るため、この層を持つ意味がある。

### `inbox/`

未整理メモや取り込み途中の素材を置く。一時置き場であり、正本ではない。

### `attachments/`

元資料を必要に応じて置く。古い職務経歴書、面接メモ、自己紹介文などを保存できるが、AI が無条件に全件参照する前提にはしない。

### `meta/`

workspace 設定、schema version、運用メモなどを置く。ユーザーキャリアの本体ではないが、workspace を維持するための補助情報を持つ。

## Layout 原則

- facts / profiles / exports を分離する
- 人間が見て意味のわかる名前にする
- AI の編集対象が directory 単位で説明できるようにする
- 一時ファイルと正本を混在させない

## AI から見た意味

workspace layout は単なる整理ではなく、AI の能力境界でもある。たとえば `facts/` への更新は慎重に扱い、`exports/` は再生成可能な派生物として扱う。どの directory に何があるかが曖昧だと、AI も人間も正本を見失いやすい。

## 未確定事項

- ファイル形式を markdown 中心にするか YAML / JSON を併用するかは未確定
- `attachments/` を MVP で必須にするかは未確定
- `embeddings/` を workspace 配下に置くかどうかは未確定

MVP では embeddings を前提機能にしない。

- embeddings は未使用でも成立する設計を保つ
- embedding index を hidden truth source にしない
- 将来追加する場合も canonical facts の代替にしない
