# CareerNess Implementation Architecture

この文書は、`implementation/` 配下の repository structure を architecture boundary として固定するためのものである。目的は説明資料を増やすことではなく、実装時に「どこへ何を置いてよいか」「どこへ何を置いてはいけないか」を明示し、human developer と AI coding agent の両方に同じ制約を共有させることにある。

## Core Principles

- Local-first を前提にする
- Workspace が source of truth である
- AI は proposal するが truth owner ではない
- Patch proposal model を中心に据える
- Structured career data を facts / profiles / exports に分離する
- Convenience より boundary を優先する

## Repository Structure Overview

```text
implementation/
├── apps/
│   ├── web/
│   └── api/
├── packages/
│   ├── workspace-core/
│   ├── patch-engine/
│   ├── schema/
│   └── shared/
├── infra/
│   ├── docker/
│   ├── terraform/
│   └── aws/
├── scripts/
├── examples/
├── .github/
│   └── workflows/
└── tmp/
```

### 意味

- `apps/` は実行可能な entrypoint を置く
- `packages/` は CareerNess の core logic と domain model を置く
- `infra/` は deploy / environment definition を置く
- `scripts/` は開発補助を置く
- `examples/` は sample workspace / patch / learning material を置く
- `.github/` は repository automation を置く
- `tmp/` は一時生成物の置き場であり、正本を置かない

## Dependency Direction

基本方向は次の通りである。

```text
apps/web
  ↓
apps/api
  ↓
packages/*

infra/*, scripts, examples, .github
  ↓
packages/* と apps/* を参照してよい
```

より厳密には、次の制約を守る。

- `apps/*` は `packages/*` に依存してよい
- `packages/*` は `apps/*` に依存してはならない
- `packages/shared` は薄い共通 utility のみに留める
- `packages/workspace-core` は workspace truth ownership を持つ
- `packages/patch-engine` は semantic patch model を持つ
- `packages/schema` は data contract を持つ
- `infra/*` は deploy concern を持つが domain truth は持たない
- `scripts/` は developer convenience のみを持ち、本番責務を持たない

## Why `apps/` と `packages/` を分離するか

CareerNess の中核は「UI 付きの AI チャット」ではない。中核は次の 4 点である。

- Workspace model
- Fact model
- Patch model
- Capability boundary

このため、entrypoint と core domain logic を分離する必要がある。

`apps/` に domain logic を寄せ始めると、次の drift が起きる。

- `apps/api` が hidden truth owner になる
- `apps/web` が canonical state を持ち始める
- AI orchestration convenience のために patch review を省略しやすくなる
- Workspace ownership が薄れ、local-first が形骸化する

`packages/` に core を固定することで、UI や orchestration を入れ替えても CareerNess の本質を維持しやすくする。

## Truth Ownership

CareerNess が守る truth ownership は明確である。

- CareerVault workspace が canonical source である
- facts は profile や export より上位の正本に近い
- profile は derived view であり truth ではない
- export は提出形式であり truth ではない
- AppServer は session と orchestration を持つが、career truth は持たない
- Browser は UI state を持つが、career truth は持たない

この ownership を repository structure に反映すると、`workspace-core` と `schema` が重要になる。

## Local-first Architecture

Local-first は「クラウドを使わない」という意味ではない。意味は次の通りである。

- ユーザーの career data の正本は local workspace に置く
- cloud は認証、AI access、session mediation の補助に留める
- apply は workspace boundary の内側で起こる
- attach していない領域を暗黙に truth source にしない

このため、`apps/api` や `infra/` に persistence convenience を集める設計は避ける。

## Workspace Ownership

Workspace ownership は `packages/workspace-core` に固定する。

理由:

- patch apply と validation を UI や API から切り離すため
- rollback / history / revision check を workspace concern に保つため
- unrestricted file access を禁止し、workspace-scoped access に限定するため

`apps/api` は patch proposal を生成してよいが、workspace truth を所有してはならない。

## AI Orchestration Boundary

AI は強いが、CareerNess の owner ではない。AI orchestration boundary は次の通りである。

- AI orchestration は `apps/api` 側に置く
- AI が生成するのは patch proposal / draft / suggestion である
- apply decision は human approval と workspace validation を通す
- AI response text をそのまま正本ファイルへ write してはならない

つまり、AI convenience のために `workspace-core` に orchestration を混ぜてはいけない。

## Package Boundary

### `packages/workspace-core`

- Workspace access
- YAML / structured file management
- validation
- apply / rollback / history

Must Not:

- AI orchestration
- cloud persistence
- UI state
- unrestricted filesystem scan

### `packages/patch-engine`

- semantic patch model
- patch diff / validation helper
- atomic change representation

Must Not:

- workspace truth ownership
- direct cloud session handling
- UI-specific patch rendering concern

### `packages/schema`

- fact / profile / export schema
- structured career data contract
- validation rules shared across runtime

Must Not:

- AI prompting logic
- workspace mutation orchestration
- presentation-only wording ownership

### `packages/shared`

- minimal shared types / utility

Must Not:

- dumping ground 化
- domain logic の逃がし先
- cross-layer shortcut dependency

## Directory Rules

### `apps/`

- 動作する product surface を置く
- domain truth を持たない

### `infra/`

- deploy と environment を表現する
- product truth や runtime business logic を持たない

### `scripts/`

- 開発作業の補助に徹する
- 本番処理の唯一実装を置かない

### `examples/`

- sample data と reference flow を置く
- test fixture と truth source を混同しない

### `.github/`

- CI/CD と repository automation を置く
- production authority や危険な自動 apply を持たない

### `tmp/`

- 一時成果物だけを置く
- git-tracked truth を置かない

## Practical Guardrails

実装時は次を確認する。

1. 新しい責務は `apps/*` と `packages/*` のどちらに属するか
2. そのコードは workspace truth を所有してよいか
3. patch proposal と patch apply が混ざっていないか
4. fact / profile / export を横断して 1 変更に詰め込んでいないか
5. convenience のために boundary を破っていないか

## Non-goals

- enterprise 向け巨大分散アーキテクチャの先回り
- cloud sync 中心設計
- AI fully autonomous editing
- repository 内に全ての将来構想を確定させること

## Future Considerations

- package 数は増えてよいが、truth ownership を曖昧にしない
- profile generation や export pipeline の独立 package 化余地はある
- auth や session は強化されてよいが、career truth を cloud に寄せない
- AWS / Terraform は成長してよいが、MVP では lightweight を維持する

## Open TODO

- package manager / build system の詳細は未確定
- monorepo tooling の標準化は未確定
- workspace revision model の具体仕様は未確定
- patch hash / approval record format の詳細は未確定
