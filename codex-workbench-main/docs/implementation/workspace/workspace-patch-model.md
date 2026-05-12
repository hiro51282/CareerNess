# Workspace Patch Model

workspace patch model は CareerNess 実装の中心である。ここでは patch format 自体ではなく、workspace 側で patch をどう解釈し、どのような更新単位として扱うかを定める。

## Principles

- patch は workspace update request である
- patch apply は workspace logic の責務である
- patch history は Git 非依存で持てるようにする
- patch と resulting state を分ける

## Patch Lifecycle

```text
proposed
  ↓
validated
  ↓
approved
  ↓
applied
  ↓
recorded
```

必要なら次も持つ。

- `rejected`
- `stale`
- `apply_failed`

## Workspace-Side Checks

apply 直前に workspace 側で再確認する。

- target path is inside workspace
- current revision still matches
- operation types are allowed
- resulting files still pass schema validation

## Update Categories

### Fact Updates

最も慎重に扱う。status transition と evidence 更新を確認する。

### Profile Updates

派生データ更新として扱う。source fact references を持てると望ましい。

### Export Updates

再生成前提の出力更新として扱う。

### History Updates

apply 結果の記録。成功/失敗の最小記録を残す。

## Revision Model

Git を前提にしないため、workspace 側で軽い revision hint を持ってよい。

- workspace revision number
- per-file modified timestamp
- patch history cursor

どれを採るかは未確定だが、stale patch 検知は必要である。

## Apply Semantics

- patch は all-or-nothing に近づける
- ただし failure 時の record は残す
- export だけ失敗して fact 成功、のような部分適用は慎重に扱う

MVP では multi-category patch を小さく保ち、部分失敗の複雑性を下げる方がよい。

## Prohibited Patterns

- patch apply が validation を飛ばす
- AI response text をそのまま file write する
- history なしで workspace を mutate する

## Open Questions

- batch patch をどこまで許可するか
- patch hash をどう算出するか
