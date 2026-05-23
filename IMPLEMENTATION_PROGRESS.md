# CareerNess Implementation Progress

## Overview

MVP フェーズの implementation status。Provider-agnostic extraction pipeline を中心に、「conversation → patch proposal → review → apply」ループを実装。

## ✅ Completed

### 1. Documentation

- [x] **careervault-mvp-schema.md**: MVP YAML schema 定義
  - 12 コアフィールド、fact lifecycle、file structure
  - patch との相性、merge conflict 最小化戦略

- [x] **careervault-yaml-structure.md**: v2+ roadmap（参考）
  - 詳細な拡張設計

- [x] **extraction-specification.md**: Provider-agnostic extraction design
  - ExtractionProvider interface
  - Codex app-server integration pattern
  - Mock provider での test strategy

### 2. Example Workspace

- [x] **implementation/examples/CareerVault/**
  - Complete minimal workspace with:
    - facts/{experiences, achievements, skills}.yaml
    - profiles/backend-engineer.yaml
    - exports/resume-jp.md
    - projects.yaml, tags.yaml, meta.yaml
  - 実装用のリファレンス

### 3. Go Implementation

#### Extraction Layer (Provider-Agnostic)

- [x] **models.go**: Data structures
  - ExtractedFact, YAMLFact, Patch, Operation
  - ExtractionPipelineResult

- [x] **provider.go**: Service orchestration
  - ExtractionProvider interface (provider-independent)
  - ExtractionService (orchestrator)
  - generatePatches() helper

- [x] **mock_provider.go**: MockExtractionProvider
  - Testing/MVP用実装
  - SetResponse() for test overrides
  - CodexExtractionProvider reference implementation（未実装、参考）

- [x] **validator.go**: Schema validation
  - ValidateAPIResult() - JSON structure check
  - ValidateYAMLFact() - MVP schema compliance
  - validateExtractedFact() - individual fact validation

- [x] **normalizer.go**: JSON → YAML transformation
  - NormalizeToYAMLFact() - type-specific conversion
  - inferTags() - keyword-based tag generation
  - formatSourceDetail() - extraction metadata

#### Patch Application Layer

- [x] **patch/applier.go**: Patch apply logic
  - Applier.ApplyPatch() - main entry point
  - applyUpsertFact() - fact insertion/update
  - applyUpdateFactStatus() - status transition
  - File I/O with YAML serialization

#### Workspace Schema

- [x] **workspace/schema.go**: MVP schema validation
  - ValidateFact() - MVP compliance
  - Required field checks
  - Type-specific rules

#### Tests

- [x] **validator_test.go**: Comprehensive validator tests
  - Valid/invalid API results
  - Schema validation edge cases
  - Status lifecycle

- [x] **normalizer_test.go**: Normalization tests
  - JSON → YAML conversion
  - ID generation
  - Tag inference
  - Period/detail extraction

- [x] **applier_test.go**: Patch application tests
  - File creation
  - Fact insertion/update
  - YAML persistence

### 4. Dependencies

- [x] **go.mod**: Cleaned up
  - Only yaml.v3 (no vendor lock-in)
  - Anthropic SDK removed
  - Ready for CodexExtractionProvider integration

### 5. API Endpoint Integration

- [x] **POST /api/v1/extract endpoint**
  - Request: `{conversation, session_id}`
  - Response: `{patches, extraction_quality, yaml_facts}`
  - ExtractionService integration
  - Mock provider for MVP
  - Error handling (empty conversation, extraction failure)

- [x] **POST /api/v1/apply-patch endpoint**
  - Request: `{patch, workspace_path}`
  - Path traversal prevention
  - applier.ApplyPatch() call
  - Response: `{patch_id, applied_count, failed_ops, updated_facts, applied_at}`

- [x] **Handler Files**
  - `internal/handler/extract.go` - PostExtract
  - `internal/handler/apply.go` - PostApplyPatch
  - Both registered in cmd/server/main.go

- [x] **JSON Serialization**
  - Added json tags to Patch, Operation, YAMLFact
  - Added json tags to ExtractionPipelineResult
  - Added json tags to ApplyResult

## 🟡 In Progress / Pending

### UI/Review Layer

- [ ] Patch review UI (mock for MVP)
  - Before/after YAML diff viewer
  - Approve/reject buttons
  - Clarification questions display

### Integration Tests

- [ ] End-to-end extraction pipeline
  - Mock provider → validation → patch generation
  - applier integration
  - File state verification

## 🔴 Not Started / v1.1+

- [ ] CodexExtractionProvider implementation
  - Codex app-server integration
  - User's LLM provider routing
  - Real fact extraction

- [ ] Frontend (web) implementation
  - Chat interface
  - Workspace management
  - Patch review UI

- [ ] Sync / Auth
  - Multi-device support
  - Session management
  - User authentication

## Key Design Decisions

### 1. Provider-Agnostic Extraction

```
ExtractionProvider interface
    ├─ MockExtractionProvider (MVP)
    ├─ CodexExtractionProvider (production - via Codex app-server)
    └─ (Other providers: OpenAI, Anthropic direct, etc.)
```

**Why**: 
- Backend は LLM vendor に依存しない
- User の own AI credits を使用
- Codex app-server が provider abstraction を担当
- Test が mock provider で可能

### 2. MVP Scope

Focus: **conversation → extraction → validation → patch → review → apply**

NOT focus:
- Advanced AI inference
- Semantic deduplication
- Graph relationships
- Cloud sync

### 3. Patch-Oriented

すべての fact 追加 / status update は patch proposal 経由
- User approval を明示的に
- Rollback が granular
- Validation が多重層

### 4. YAML as Source of Truth

- facts/*.yaml が facts の canonical source
- Local-first operation
- Git-friendly diffs
- No hidden state in DB

## Testing Status

✅ **Unit Tests**:
- validator_test.go: PASS
- normalizer_test.go: PASS

⚠️  **Integration Tests**:
- applier_test.go: Ready (no SDK dependency)
- End-to-end: Pending (needs API integration)

## Next Steps (Priority Order)

1. ✅ **API Endpoint Integration** (COMPLETED)
   - Wire ExtractionService to HTTP handler ✅
   - Request validation ✅
   - Response formatting ✅

2. **End-to-End Integration Test** (2-3 hours)
   - extraction → validation → patch → applier full loop
   - File system verification with test workspace
   - Error case testing

3. **Patch Review UI** (4-6 hours)
   - Frontend: before/after YAML diff viewer
   - Approve/reject interaction
   - Mock patch display with quality metrics

4. **CodexExtractionProvider** (next phase)
   - Not required for MVP
   - Reference implementation provided in mock_provider.go

## Technical Debt / Known Issues

- [ ] go.sum cleanup (ensure yaml.v3 only)
- [ ] provider.go: generatePatchFromFact() duplication check
- [ ] Mock provider: keyword-based extraction は obviously not production
- [ ] Applier: update_fact_status operation parsing needs refinement
- [ ] Tests: applier_test integration with extraction service

## Files Summary

```
implementation/apps/api/
├── internal/
│   ├── extraction/
│   │   ├── models.go (data structures)
│   │   ├── provider.go (service + interface)
│   │   ├── mock_provider.go (MVP provider)
│   │   ├── validator.go (validation logic)
│   │   ├── normalizer.go (transformation)
│   │   ├── validator_test.go ✅
│   │   ├── normalizer_test.go ✅
│   │   └── README.md (architecture)
│   ├── patch/
│   │   ├── applier.go (patch apply logic)
│   │   └── applier_test.go (ready for testing)
│   └── workspace/
│       └── schema.go (MVP schema validation)
├── go.mod (cleaned - yaml.v3 only)
└── go.sum
```

## Deployment Readiness

**MVP**: ✅ Ready to implement API endpoints
- Core extraction logic: complete
- Validation: complete
- Patch generation: complete
- Patch application: complete
- Tests: passing (validator, normalizer)

**v1.0+**: Need Codex app-server integration
- CodexExtractionProvider: not needed for MVP
- Frontend: separate phase

## Summary

Provider-agnostic extraction pipeline + deterministic patch workflow が実装完了。Backend は YAML schema handling と validation に集中。LLM provider は CodexExtractionProvider 経由で交換可能。

**MVP での deliverable**: conversation → patch proposal loop（LLM provider mock）
