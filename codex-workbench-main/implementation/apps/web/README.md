# apps/web

`apps/web/` は CareerNess の React frontend を置く。役割は人間の確認、比較、承認、編集体験を提供することであり、career data の truth owner になることではない。

## Responsibility

- UI rendering
- session-scoped state の管理
- patch review UI
- profile compare UI
- workspace attach / conversation / approval の操作導線

## Must Not

- facts / profiles / exports の canonical source にならない
- workspace apply logic を独自実装しない
- attach していない local data を暗黙に探索しない
- AI proposal を承認なしで truth として確定しない

## Why This Exists

- 人間が patch proposal を理解し、レビュー可能な形で意思決定するため
- Local-first でも UX を犠牲にせず、change visibility を高めるため

## Owned Concepts

- view state
- session UI state
- patch review interaction
- compare / diff presentation
- approval affordance

## Dependencies

- `packages/schema`
- `packages/patch-engine`
- API contract 経由の `apps/api`

## Must Not Depend On

- `packages/workspace-core` の内部ファイル操作実装
- cloud persistence を前提にした hidden truth
- backend 内部都合に密結合した UI state

## Future Considerations

- offline-friendly UI は強化してよい
- richer diff preview や validation hint は追加してよい
- ただし review と apply を一体化しすぎない

## Non-goals

- workspace ownership
- AI orchestration ownership
- career data の長期保存
