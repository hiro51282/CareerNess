# 002 Patch-Oriented Editing

## Status

Accepted

## Context

CareerNess は AI-assisted, not AI-owned を守る必要がある。AI が直接 workspace を自由編集すると、何が変わったか、なぜ変わったか、どこまで戻せるかが不透明になる。

## Decision

workspace 更新は patch proposal model を採用する。

```text
AI
  ↓
Patch proposal
  ↓
User review
  ↓
Apply
  ↓
Workspace update
```

AI は patch proposal を生成するが、承認前の正本更新は行わない。

## Consequences

### Positive

- change visibility が高い
- rollback 単位を作りやすい
- hallucinated fact の流入を抑えやすい
- Local-first と相性がよい

### Negative

- 直接編集より操作数は増える
- diff / approval UI の実装が必要になる

## Rejected Alternatives

- 会話送信のたびに自動保存する
- profile/export 更新は承認不要としてまとめて apply する
- AI に unrestricted file write 権限を与える
