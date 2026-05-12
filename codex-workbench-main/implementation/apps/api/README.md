# apps/api

`apps/api/` は Go backend / orchestration layer を置く。ここは CareerNess の実行制御と AI coordination を担うが、workspace の正本所有者ではない。

## Responsibility

- AppServer integration
- auth / session mediation
- AI orchestration
- patch proposal generation
- Browser と workspace-scoped capability の仲介

## Must Not

- canonical workspace ownership を持たない
- hidden truth persistence を追加しない
- AI response をそのまま確定保存しない
- unrestricted filesystem access を提供しない

## Why This Exists

- frontend から auth, session, AI coordination を分離するため
- workspace attach と AI capability を調停し、無制限な権限拡大を防ぐため

## Owned Concepts

- authenticated request handling
- orchestration flow
- AI tool / model invocation boundary
- patch proposal lifecycle の server-side mediation

## Dependencies

- `packages/patch-engine`
- `packages/schema`
- `packages/workspace-core` の公開 API
- 必要最小限の `packages/shared`

## Must Not Depend On

- `apps/web` の UI 実装詳細
- workspace の hidden mirror storage
- `tmp/` や examples を本番状態として扱うこと

## Future Considerations

- session hardening や auditability は強化してよい
- orchestration は厚くなってよいが、truth ownership だけは移してはいけない

## Non-goals

- canonical database layer
- full autonomous apply engine
- career fact generation の正本化
