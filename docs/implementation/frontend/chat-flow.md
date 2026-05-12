# Chat Flow

CareerNess の chat は単なる会話 UI ではなく、fact extraction と patch review へつながる入口である。したがって、雑談チャットと同じ設計にはしない。

## Goal

- ヒアリングと構造化を支援する
- AI proposal と workspace update の境界を保つ
- clarification と approval を自然につなぐ

## Core Flow

```text
User message
  ↓
AI response
  ↓
clarification or proposal
  ↓
patch review
  ↓
apply or continue chatting
```

## Modes In One Thread

同じ会話内でも、少なくとも次の状態を意識する。

- discovery
- clarification
- proposal preview
- review decision
- post-apply summary

## UI Requirements

- 現在 attach されている workspace が見える
- AI が参照した対象範囲を説明できる
- proposal が会話文と区別される
- approval pending 状態が明確である

## Fact-Oriented Conversation

質問の重心は肩書ではなく、行動・判断・成果に寄せる。

良い質問例:

- 何を改善したか
- 何を判断したか
- どんな制約があったか
- 誰と調整したか
- どんな影響が出たか

## Prohibited UX

- 会話送信だけで保存完了に見える UI
- proposed fact と confirmed fact の見分けがつかない UI
- profile draft と export draft を fact 同等に見せる UI

## Open Questions

- timeline 型 UI を入れるか
- review panel を side-by-side にするか
