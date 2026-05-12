# Frontend Structure

frontend は CareerVault の正本を持たないが、ユーザーが boundary と diff を理解するための重要な層である。派手な stateful editor より、責務が見える構造を優先する。

## Main Areas

- session / auth
- workspace attach
- chat / extraction
- patch review
- fact / profile / export browsing

## Suggested Structure

### Session Shell

ログイン状態、workspace attach 状態、現在の会話や review 状態を保持する。

### Workspace Explorer

facts / profiles / exports / projects を読み分けて見せる。これは正本エディタではなく、閲覧と確認の入口である。

### Chat Workspace

AI と会話し、clarification や proposal preview を扱う。

### Review Workspace

patch diff、rationale、risk を見せる。

### Generation Views

profile preview と export preview を表示する。

## State Separation

frontend では次を分ける。

- conversational state
- attachment state
- proposal state
- approved/applied state

proposal state を通常の chat state に埋め込むだけにすると、承認と保存の意味がぼやける。

## Responsibility Limits

- canonical patch apply logic を持ちすぎない
- fact truth をブラウザ内 only state にしない
- hidden auto-save を正本扱いしない

## Open Questions

- local file access helper との境界をどう見せるか
