# Implementation Docs

このディレクトリは、CareerNess の実装設計書を置く場所である。目的は framework 選定メモを書くことではなく、人間と AI coding agent が同じ制約を共有したまま実装できるようにすることにある。

## Core Invariants

- Local-first
- User-owned data
- Workspace-scoped AI
- Patch proposal model
- Structured career data
- 「肩書」ではなく「行動・判断・成果」
- AI-assisted, not AI-owned
- Human-readable design

## Reading Order

1. `decisions/`
2. `workspace/`
3. `ai/`
4. `profile/` and `export/`
5. `backend/`, `frontend/`, `auth/`, `deploy/`

## Most Important Documents

- `workspace/fact-schema.md`
- `workspace/workspace-patch-model.md`
- `workspace/history-and-rollback.md`
- `ai/ai-patch-format.md`
- `ai/approval-flow.md`
- `ai/tooling-model.md`
- `profile/profile-generation.md`
- `export/export-pipeline.md`

## Implementation Attitude

- AppServer は orchestration layer であり truth owner ではない
- CareerVault workspace が source of truth である
- facts / profiles / exports を混同しない
- proposal と apply を分離する
- AI convenience より change visibility と user trust を優先する
