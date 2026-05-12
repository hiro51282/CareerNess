# Workspace Patch Model

workspace patch model は CareerNess 実装の中心である。ここでは patch format 自体ではなく、workspace 側で patch をどう解釈し、どのような更新単位として扱うかを定める。

## Principles

- patch は workspace update request である
- patch 1 つは 1 semantic change を表す
- patch apply は workspace logic の責務である
- patch history は Git 非依存で持てるようにする
- patch と resulting state を分ける
- AppServer は patch proposal を運ぶが truth owner にはならない

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

## Semantic Atomicity

CareerNess の patch model では、reviewability を優先する。

- 1 patch = 1 semantic change を原則とする
- fact 追加と profile 再生成は別 patch に分ける
- profile 再生成と export 再生成も別 patch に分ける
- 複数 file を触っても、意味上 1 変更なら 1 patch でよい

悪い例:

- facts 追加、status 昇格、profile 再生成、export 更新を 1 patch に詰め込む

良い例:

- 1 patch で 1 fact を追加する
- 別 patch で profile draft を更新する
- 別 patch で export draft を更新する

## Workspace-Side Checks

apply 直前に workspace 側で再確認する。

- target path is inside workspace
- current revision still matches
- operation types are allowed
- resulting files still pass schema validation
- approval scope がその semantic change に対応している

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

MVP では semantic atomicity を保ち、multi-category patch を避けることで部分失敗の複雑性を下げる方がよい。

## Prohibited Patterns

- patch apply が validation を飛ばす
- AI response text をそのまま file write する
- history なしで workspace を mutate する
- review 不可能な巨大 patch を 1 回で apply する

## Open Questions

- batch patch をどこまで許可するか
- patch hash をどう算出するか
