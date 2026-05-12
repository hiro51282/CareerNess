# API Flow

backend API は Browser と AppServer の間で、会話、workspace attach、patch proposal、approval apply をつなぐ層である。API は domain truth を所有しないが、責務境界を壊さないための入口として重要である。

## Goal

- 会話系操作と workspace 操作を分けて表現する
- proposal と apply を別 endpoint / action にする
- Local-first と patch proposal model を API でも崩さない

## Core Flows

### Conversation to Proposal

```text
POST conversation message
  ↓
AppServer reads attached workspace context
  ↓
AI extraction / synthesis
  ↓
Patch proposal returned
```

### Approval to Apply

```text
POST approval decision
  ↓
Patch re-validation
  ↓
Local runtime / workspace apply
  ↓
History append
  ↓
Updated workspace summary returned
```

## API Responsibility Split

### Conversation APIs

- ユーザー発話送信
- AI 応答取得
- clarification question 取得
- proposal preview 取得

### Workspace APIs

- workspace attach / detach
- directory summary 読み取り
- file preview
- patch history 取得

### Patch APIs

- patch proposal 生成
- patch validation
- approval submit
- approved patch apply
- rollback patch prepare

## Recommended Resource Shape

細かい URI 設計は固定しないが、責務は分ける。

- `/sessions`
- `/workspaces`
- `/conversations`
- `/patches`
- `/approvals`
- `/history`

## API Rules

- `generate proposal` と `apply patch` を同じ call にしない
- `approve patch` と `apply patch` も概念上は分離する
- fact 更新 API は status transition を明示する
- export regenerate は fact write API を兼ねない
- workspace root は session-scoped attachment から解決する
- session-only inference を hidden write として扱わない

## Validation Layer

API では少なくとも次を担保する。

- attach 済み workspace か
- patch target が workspace 配下か
- approval state が存在するか
- patch hash または revision が一致するか

## Streaming Considerations

AI 応答は streaming してよいが、patch proposal は最終確定ブロックとして分ける。

- conversational text
- structured proposal

この 2 層を混ぜると、UI 側で「返答」と「変更案」の境界が曖昧になる。

## Prohibited Patterns

- chat response の副作用として patch を自動 apply する
- approval endpoint が raw file path を自由入力で受ける
- history append を失敗しても成功扱いにする
- 1 API call で複数 semantic change を不可分に apply する

## Open Questions

- SSE と WebSocket のどちらを主にするかは未確定
- local helper process を API に含めるかは未確定
