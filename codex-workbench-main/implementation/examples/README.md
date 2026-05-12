# examples

`examples/` は sample workspace、sample patch、reference flow などを置く。AI coding agent の learning material としても使えるが、truth source ではない。

## Responsibility

- example CareerVault
- sample patch proposal
- reference data flow
- onboarding / documentation support

## Must Not

- 本番 canonical data を置かない
- 実装の暗黙仕様を example のみで定義しない
- test fixture と product requirement を混同しない

## Why This Exists

- 人間と AI の両方が期待構造を素早く理解できるようにするため
- boundary-heavy architecture を具体例で示すため

## Owned Concepts

- illustrative sample
- reference fixture
- educational artifact

## Dependencies

- `packages/schema`
- `packages/patch-engine`
- `docs/` の設計文脈

## Must Not Depend On

- hidden implementation detail
- runtime の偶発的挙動

## Future Considerations

- example は最小でも高品質に保つ
- drift を防ぐため、実際の schema / patch model と同期させる

## Non-goals

- integration test の唯一基盤
- product behavior の正本
