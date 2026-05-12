# 001 Workspace Boundary

## Status

Accepted

## Context

CareerNess は user-owned career data を扱う。利便性のために AI へ広いローカル権限を渡すと、Workspace-scoped AI という前提が崩れ、設計の説明責任も失われる。

## Decision

AI が参照・提案対象にできるデータは、ユーザーが明示的に attach した workspace に限定する。

- workspace root を capability boundary とする
- attach されていないパスは参照対象にしない
- patch target も workspace 配下に限定する
- AppServer は workspace 外の補完知識を正本候補として持ち込まない

## Consequences

### Positive

- ownership が明確になる
- AI の能力説明が簡潔になる
- review 対象が workspace 内に閉じる

### Negative

- 既存レジュメやメモを自動で拾う convenience は下がる
- 初回 attach と整理の UX が必要になる

## Rejected Alternatives

- ホームディレクトリ全体を AI 検索対象にする
- AppServer 側にユーザーデータを継続収集して補完する
- export 済み文章から facts を暗黙復元する
