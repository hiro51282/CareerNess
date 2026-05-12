# .github/workflows

`implementation/.github/workflows/` は CI/CD workflows を置く。ここは品質ゲートと補助自動化の置き場であり、危険な authority を持たせない。

## Responsibility

- CI validation
- build / test / lint automation
- deploy pipeline の最小定義

## Must Not

- AI-generated code や patch を無審査で本番適用しない
- workspace truth を workflow 内に長期保持しない
- repository 外の高権限操作を暗黙に増やさない

## Why This Exists

- 継続的検証を自動化しつつ、人間の承認境界を守るため

## Owned Concepts

- verification gate
- release automation
- safe repository event handling

## Dependencies

- `scripts/`
- build/test 対象の `apps/*`, `packages/*`
- deploy 時は `infra/*`

## Must Not Depend On

- manual-only knowledge
- tmp の偶発的ファイル
- hidden production truth

## Future Considerations

- preview deploy や security scan の追加余地がある
- ただし workflow 自体が product owner になってはいけない

## Non-goals

- autonomous operations platform
- patch approval bypass
