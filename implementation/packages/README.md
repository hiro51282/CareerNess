# packages

`packages/` は CareerNess の core domain logic を置く。CareerNess の本質は `apps/` ではなく、この層に寄せる。

## Responsibility

- workspace model
- patch model
- schema / validation
- capability boundary を支える core logic

## Must Not

- app entrypoint concern を持たない
- UI state や auth session を core domain に混ぜない
- deploy/environment concern を抱え込まない

## Why This Exists

- CareerNess を UI や server runtime から独立した core system として保つため
- AI coding agent に「ここが本質ロジック」と明示するため

## Owned Concepts

- truth ownership rules
- semantic patch rules
- structured career data rules
- workspace-scoped mutation rules

## Dependencies

- package 間の依存は明示的かつ一方向に保つ
- 外部 library は package responsibility に必要な範囲に限定する

## Must Not Depend On

- `apps/*`
- `infra/*`
- `.github/*`
- 本番以外の script convenience

## Future Considerations

- package の追加は許容する
- ただし「なんとなく shared」な package 増殖は避ける

## Non-goals

- monorepo の雑多な共通置き場
- framework 設定の収容場所
