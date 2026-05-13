# packages/patch-engine

`packages/patch-engine/` は patch proposal / diff / apply model を扱う。CareerNess の patch-oriented editing を成り立たせる中核 package である。

## Responsibility

- patch proposal model
- diff representation
- patch validation helper
- patch atomicity rule
- 1 patch = 1 semantic change の原則を支える共通ロジック

## Must Not

- workspace truth ownership を持たない
- direct file write の最終責務を持たない
- AI orchestration を持たない
- UI 固有の diff rendering concern を内包しない

## Why This Exists

- proposal と apply を分離し、reviewability を上げるため
- semantic change を transport や UI 実装から独立させるため

## Owned Concepts

- patch schema
- semantic atomicity
- diff unit
- patch status / lifecycle metadata

## Dependencies

- `packages/schema`
- 必要最小限の `packages/shared`

## Must Not Depend On

- `apps/api`
- `apps/web`
- cloud session / auth logic
- workspace-specific hidden state

## Future Considerations

- richer patch metadata や conflict hint は追加してよい
- batch patch は将来検討してよい
- ただし巨大 patch の常態化は避ける

## Non-goals

- canonical workspace storage
- full UI diff system
- unrestricted code-mod engine
