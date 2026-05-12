# History And Rollback

CareerNess は Git を前提にしない。ただし、戻せない local-first system は user-owned data に対して不誠実である。そこで、軽量な history と rollback philosophy を持つ。

## Principles

- Git は optional
- patch history は product responsibility として持つ
- full version control system を再実装しない
- 「何が変わったか」と「どこまで戻せるか」を最低限残す

## What To Record

- patch id
- created_at / applied_at
- patch summary
- affected files
- approval result
- apply result
- rollback preparation info

## Why History Matters

- AI proposal の誤りを戻せる
- ユーザー manual edit と AI apply の境界を追える
- facts 更新の由来を説明しやすい

## Rollback Model

MVP では「任意時点への完全復元」までは目指さない。最低限、直近 patch を安全に戻せる設計を優先する。

### Option A

inverse patch を記録する。

### Option B

変更対象ファイルの pre-apply snapshot を軽く残す。

どちらを採るかは未確定だが、Git 非前提で説明できる必要がある。

## Scope

- facts rollback は重要
- profiles rollback はあると望ましい
- exports rollback は優先度低め

## Prohibited Patterns

- apply history を残さない
- rollback 不可能な destructive delete を標準操作にする
- history が cloud-only で local workspace に残らない

## Open Questions

- history ディレクトリ構造
- snapshot retention policy
- rejected patch をどこまで残すか
