# Session Model

CareerNess では session を一つに潰さない。少なくとも user session、workspace attachment、conversation state を分けて扱う必要がある。

## Session Types

### User Session

認証済みユーザーとして AppServer を利用する単位。

### Workspace Attachment

どの local workspace にアクセスしてよいかを表す単位。user session より短命でもよい。

### Conversation Session

ヒアリングや profile 生成など、会話の文脈を持つ単位。

### Patch Review Session

生成済み patch proposal を review / approval する単位。conversation と 1:1 でなくてよい。

## Why Separate Them

- ログインしていても workspace 未接続はありうる
- 同じ workspace で複数会話はありうる
- 会話が続いていても patch review を一時中断したい
- approval 後 apply 前に再 validation が必要な場合がある

## Minimal State

### User Session

- user id
- auth status
- AI account availability summary

### Workspace Attachment

- workspace root identifier
- attachment status
- allowed capability scope
- current workspace revision hint

### Conversation Session

- conversation id
- selected task mode
- referenced workspace slices

### Patch Review Session

- patch id
- approval state
- review hash or revision

## Lifecycle

```text
login
  ↓
user session created
  ↓
workspace attached
  ↓
conversation started
  ↓
patch proposed
  ↓
review session opened
  ↓
approved / rejected
  ↓
apply or discard
```

## Expiry Behavior

- user session expiry は認証レイヤ依存
- workspace attachment expiry は短めでよい
- conversation session は resume 可能でもよい
- patch review session は stale 判定が必要

## Staleness Handling

patch proposal 作成後に workspace が変わった場合、apply 前に stale 判定する。

- target file changed
- target fact changed
- patch history advanced

stale の場合は silent rebase ではなく、再 validation または再承認を促す。

## Prohibited Patterns

- conversation id だけで apply を許可する
- stale patch を自動書き換えして apply する
- expired attachment でも workspace read を継続する

## Open Questions

- session persistence をどこまで local resume するかは未確定
