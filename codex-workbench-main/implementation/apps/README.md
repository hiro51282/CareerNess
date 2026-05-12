# apps

`apps/` は CareerNess の実行可能な application entrypoint 群を置くディレクトリである。ここは product surface の層であり、core domain truth の置き場ではない。

## Responsibility

- 実行可能な Web / API application を配置する
- user interaction, session, request/response の入口を提供する
- `packages/` の core logic を利用して product behavior を組み立てる

## Must Not

- domain logic の正本を持たない
- workspace canonical logic を `apps/` 配下に重複実装しない
- patch apply の最終責務を持たない
- convenience のために `apps/api` を hidden data owner にしない

## Why This Exists

- entrypoint と core packages を分離し、UI や backend orchestration の変更が truth ownership を壊さないようにするため
- AI coding agent に「動作面」と「本質ロジック」を見分けさせるため

## Owned Concepts

- application bootstrapping
- request lifecycle
- session lifecycle
- UI/API surface

## Dependencies

- `packages/workspace-core`
- `packages/patch-engine`
- `packages/schema`
- 必要最小限の `packages/shared`

## Must Not Depend On

- `apps/*` 同士の密結合な内部実装共有
- `tmp/` の生成物を前提にした本番挙動
- examples を本番 truth とみなす実装

## Future Considerations

- app が増えても domain logic は `packages/` に寄せる
- desktop / CLI app を追加しても同じ boundary を維持する

## Non-goals

- monolith 的に全責務を抱えること
- backend convenience のために architecture の中心になること
