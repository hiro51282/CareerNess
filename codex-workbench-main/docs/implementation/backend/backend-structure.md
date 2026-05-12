# Backend Structure

CareerNess の backend は AppServer であり、AI orchestration layer としての責務を持つ。複雑な layered architecture を先に固定するより、壊してはいけない境界を先に固定する。

## Core Responsibility

- user session を扱う
- workspace attachment を扱う
- AI conversation を扱う
- patch proposal / validation / apply orchestration を扱う

## Must Not

- CareerVault の source of truth になる
- facts / profiles / exports を恒久保存する
- unrestricted local access の代理になる

## Suggested Modules

### Session Layer

- auth session
- workspace attachment state
- conversation session state

### AI Orchestration Layer

- prompt assembly
- workspace context selection
- extraction / synthesis dispatch
- tool invocation control

### Patch Layer

- patch proposal builder
- schema validation
- risk classification
- approval token check

### Workspace Gateway

- read-only listing / preview
- approved apply
- history append
- rollback patch generation

### Observability Layer

- operational logs
- error tracing
- minimal audit metadata

## Data Handling Rule

backend に残るのは transient context を原則とする。

- keep: session metadata, request correlation, minimal operational logs
- avoid: full workspace mirror, long-lived career data cache

## Boundary With Workspace

workspace structure と canonical schema は local workspace 側に寄せる。backend はそれを解釈し補助するが、独自 canonical schema を内側に持って上書きしてはいけない。

## Boundary With Frontend

frontend は閲覧と操作の入口であり、backend は session / AI / patch orchestration を担う。frontend が canonical patch logic を持ちすぎると再利用しづらくなり、backend が UI state を持ちすぎると責務が濁る。

## Implementation Bias

MVP では小さく保つ。

- thin API surface
- explicit patch objects
- simple workspace gateway
- minimal background jobs

## Prohibited Patterns

- backend 内部 DB に canonical facts を複製する
- profile/export の再生成結果だけを backend に残し workspace に戻さない
- AI prompt 用 convenience cache を truth 扱いする

## Open Questions

- local helper を backend module とみなすか別 runtime とみなすかは未確定
