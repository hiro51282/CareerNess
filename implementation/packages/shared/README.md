# packages/shared

`packages/shared/` は薄い shared utility と type を置く。便利だから何でも入れる場所ではない。

## Responsibility

- cross-package で本当に再利用される最小の type
- domain 非依存または低依存の utility
- boundary を壊さない共通定義

## Must Not

- misc dumping ground にならない
- domain logic の避難所にならない
- `workspace-core` や `patch-engine` の責務を抜き出して薄めない
- app-specific helper を共有資産扱いしない

## Why This Exists

- package 間で同じ低レベル定義を重複させないため
- ただし shared 化のコストを意識し、責務の所在を曖昧にしないため

## Owned Concepts

- primitive shared types
- low-level utility
- cross-package constants のうち本当に中立なもの

## Dependencies

- できるだけ依存を持たない
- 持つとしても軽量かつ domain-neutral なものに限る

## Must Not Depend On

- `packages/workspace-core`
- `packages/patch-engine`
- `packages/schema`
- `apps/*`

## Future Considerations

- shared が肥大化したら再分割する
- 中立に見えても domain を持ち始めたら元 package へ戻す

## Non-goals

- 共通化それ自体
- 責務の判断を先送りすること
