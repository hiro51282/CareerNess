# 003 Local-First Storage

## Status

Accepted

## Context

CareerNess は career information の ownership をユーザー側に置く。クラウド中心ストレージは便利だが、正本の所在と retention が曖昧になりやすい。

## Decision

CareerVault を local workspace として扱い、facts / profiles / exports の正本および派生物は原則ここに保持する。

- YAML ベースを優先する
- human-readable / diff-friendly を優先する
- AppServer は orchestration layer にとどめる
- Git は前提にしない

## Consequences

### Positive

- user-owned data を守りやすい
- backup / file sync の選択をユーザーに委ねられる
- AI と人間の双方にとって可視な正本になる

### Negative

- デバイス喪失時の保護はユーザー責任が増える
- multi-device sync は別設計が必要になる

## Rejected Alternatives

- AppServer 内 DB を canonical store にする
- export を唯一の保存対象にする
- opaque internal database を正本にして human-readable file を副生成物にする
