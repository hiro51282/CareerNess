# State Management

CareerNess の frontend state は「何を表示しているか」だけでなく、「まだ proposal なのか、もう apply 済みなのか」を区別できる必要がある。

## State Domains

- auth state
- workspace attachment state
- conversation state
- proposal state
- review state
- applied workspace snapshot summary

## Principles

- proposal と applied result を混ぜない
- optimistic update で facts を既成事実化しない
- local draft edit と canonical workspace state を区別する

## Recommended Handling

### Auth State

ログイン状態と AI 利用可能性だけを持つ。

### Attachment State

workspace root, attachment status, revision hint を持つ。

### Conversation State

メッセージと会話モードを持つ。

### Proposal State

current patch, risk summary, validation result を持つ。

### Review State

approval decision, user edits, stale warning を持つ。

### Applied Summary State

最後に反映された patch id と更新サマリを持つ。

## Anti-Patterns

- proposal を自動で fact list に混ぜる
- export preview を save 済みに見せる
- stale patch を無視して local state だけ更新する

## Open Questions

- local cache persistence をどこまで持つかは未確定
