# packages/workspace-core

`packages/workspace-core/` は CareerVault workspace の canonical access layer である。CareerNess の中で最も重要な package の 1 つであり、workspace の source of truth を守る責務を持つ。

## Responsibility

- CareerVault access
- YAML / structured file management
- workspace-scoped file operation
- patch apply
- validation
- revision check
- rollback support
- apply / history の最小記録

## Must Not

- AI orchestration を持たない
- profile generation ownership を持たない
- cloud persistence を持たない
- unrestricted filesystem access を持たない
- UI state を持たない

## Why This Exists

- workspace を AppServer や UI から守り、truth ownership を固定するため
- patch apply を convenience write ではなく、validation 付き mutation として扱うため

## Owned Concepts

- workspace boundary
- workspace revision / stale detection
- canonical file mutation
- rollback / recovery
- apply safety check

## Dependencies

- `packages/schema`
- `packages/patch-engine`
- 必要最小限の `packages/shared`

## Must Not Depend On

- `apps/api`
- `apps/web`
- AI provider SDK や prompting logic
- cloud DB / object storage client

## Future Considerations

- workspace lock や concurrent edit 制御は強化してよい
- audit trail は強化してよい
- ただし orchestration を引き受けてはいけない

## Non-goals

- AI による自律編集の中心
- presentation layer
- repository 全体の汎用 FS abstraction
