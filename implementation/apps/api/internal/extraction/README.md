# Extraction Package

Provider-agnostic fact extraction pipeline. AI provider は交換可能部品として扱われます。

## Architecture

```
Conversation
    ↓
ExtractionService (orchestrator)
    ├─ ExtractionProvider (interface)
    │   ├─ MockExtractionProvider (MVP)
    │   ├─ CodexExtractionProvider (production - calls Codex app-server)
    │   ├─ OpenAIExtractionProvider (reference - direct OpenAI)
    │   └─ Other providers (Anthropic via Codex, etc.)
    │
    ├─ ValidateAPIResult() (provider-independent)
    ├─ NormalizeToYAMLFact() (provider-independent)
    ├─ ValidateYAMLFact() (MVP schema)
    └─ generatePatches() (patch proposal)
    ↓
ExtractionPipelineResult (patches ready for user review)
```

## Key Design

### Provider Abstraction

```go
type ExtractionProvider interface {
    ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error)
    Name() string
}
```

**Why**: CareerNess は LLM provider に依存しない。実装は以下のいずれでも可能：
- **MockExtractionProvider**: テスト、デモ用（MVP では実装済み）
- **CodexExtractionProvider**: 本番用。Codex app-server 経由で user の LLM を呼び出し
- **Direct providers**: 参考実装（使用非推奨 - security/ownership 問題）

### MVP フロー

```
User: "I spent 3 years at Acme Corp..."
    ↓
ExtractionService.ExtractFromConversation()
    ↓
Mock Provider returns structured JSON
    ↓
Validators (JSON structure, schema compliance)
    ↓
Normalizers (JSON → YAML)
    ↓
Patch generator
    ↓
Patch proposal (user review 用)
```

## Provider Integration

### CodexExtractionProvider (Production)

```
CareerNess Backend
    ↓
POST /v1/extract {conversation: "..."}
    ↓
Codex App-Server
    ├─ User authentication
    ├─ LLM provider routing (OpenAI/Anthropic/etc)
    ├─ Call user's LLM via user's API key
    └─ Return structured JSON
    ↓
CareerNess Backend receives provider-agnostic JSON
```

**利点**:
- ユーザー自身の AI credit を使用
- Backend が vendor lock-in されない
- 複数の LLM provider サポート可能
- セキュリティ（API key は Codex が管理）

### MockExtractionProvider (MVP/Testing)

```go
service := extraction.NewExtractionService(extraction.NewMockExtractionProvider())
result, err := service.ExtractFromConversation(ctx, "conversation...", "sess_123")
// Returns mock-populated ExtractionPipelineResult
```

## Usage

### MVP (テスト・デモ)

```go
// Mock provider で extraction pipeline をテスト
provider := extraction.NewMockExtractionProvider()
service := extraction.NewExtractionService(provider)

result, err := service.ExtractFromConversation(ctx, conversation, sessionID)
// result.Patches → user に送信
// result.YAMLFacts → 内部キャッシュ
```

### Production (未実装 - 参考)

```go
// Codex app-server 経由で extraction
provider := extraction.NewCodexExtractionProvider(
    "https://codex.app-server.com",
    userSessionToken,
)
service := extraction.NewExtractionService(provider)

result, err := service.ExtractFromConversation(ctx, conversation, sessionID)
// Codex が user の LLM を呼び出し
// CareerNess は provider-agnostic JSON を受け取り
```

## Testing

```go
// Mock provider を使って完全な pipeline をテスト
mock := extraction.NewMockExtractionProvider()
mock.SetResponse(&extraction.ExtractedFactResult{
    ExtractedFacts: []extraction.ExtractedFact{
        {
            Type: "experience",
            FactIDHint: "test-project",
            // ...
        },
    },
})

service := extraction.NewExtractionService(mock)
result, err := service.ExtractFromConversation(ctx, "test", "sess_test")
// Assertions on result.Patches, result.YAMLFacts
```

## Files

- **models.go**: Extracted* / YAMLFact / Patch struct definitions
- **provider.go**: ExtractionProvider interface + ExtractionService
- **mock_provider.go**: MockExtractionProvider (MVP)
- **validator.go**: JSON structure and YAML schema validation (provider-independent)
- **normalizer.go**: JSON → YAML transformation (provider-independent)
- **validator_test.go**, **normalizer_test.go**: Tests (provider-independent)
- **README.md**: このファイル

## MVP Principle

- **Provider は交換可能**: テスト、デモ、本番で異なる provider を使用可能
- **Core logic は provider-independent**: validation、normalization、patch generation は LLM に依存しない
- **Integration point は limited**: ExtractionProvider interface の ExtractFacts() だけ
- **Mock first**: MVP では MockExtractionProvider で十分。本番化時に CodexExtractionProvider を実装

## Future

v1.1+:
- CodexExtractionProvider 実装（app-server integration）
- Logging、metrics（provider call の成功率など）
- Retry logic、error recovery
- Rate limiting、quota management
- Multi-provider failover
