# 004 AI Provider Integration

## Status

Accepted

## Context

CareerNESS は AI 機能（fact extraction、profile 生成、export 生成）に LLM を使用する。
`docs/implementation/ai/extraction-specification.md` には `github.com/anthropic-ai/sdk-go` を使った実装例が記載されていたが、これはドキュメントの誤記であることが判明した。

プロダクトオーナーの追加前提として以下が確定している（2026-05-31）：

- ユーザー自身の OpenAI アカウント（Codex 経由）を利用する
- CareerNESS 運営は AI 利用料の請求主体にならない
- LLM 利用のための API キーおよび利用料金はプロダクト側では提供しない

これにより、3つのドキュメント間に存在した矛盾（CLAUDE.md の「Codex AppServer 経由」、extraction-specification.md の Anthropic SDK 実装例、auth-flow.md の「ユーザー自身の OpenAI / ChatGPT アカウント」）が解消された。

## Decision

CareerNESS の AI 呼び出しは、**ユーザー自身の OpenAI アカウント（Codex 経由）** を使う。

- AI provider: OpenAI（Codex API）
- API key: ユーザー自身が保持・提供する
- 利用料金: ユーザー負担。CareerNESS 運営は請求主体にならない
- `apps/api`（AppServer）はブラウザとユーザーの OpenAI アカウントの中継層として機能する

`extraction-specification.md` の Anthropic SDK 実装例（`claude-opus-4-7` モデルの使用）は誤記であり、OpenAI SDK を使う実装に差し替える必要がある。

## Consequences

### Positive

- CareerNESS 運営が AI 利用コストを負担しない
- ユーザーが自分の OpenAI プランの範囲でコストを管理できる
- Anthropic との契約・運用が不要
- rate limiting の目的が「OpenAI コスト管理」ではなく「AppServer abuse 防止」に限定され、設計がシンプルになる

### Negative

- ユーザーが OpenAI アカウントを持っていることが前提条件になる
- ユーザーの API key を AppServer がセキュアに扱う仕組みが必要（key の一時受け取り方式と、ブラウザから直接呼ぶ方式のどちらを採用するかは未確定）
- OpenAI のサービス変更・価格改定がユーザーに直接影響する

## Rejected Alternatives

- **Anthropic Claude API を直接呼ぶ**（extraction-specification.md の誤った実装例）: CareerNESS がユーザーの Anthropic アカウントを管理する構造になり、追加前提と矛盾する
- **CareerNESS 運営が API キーを一括管理・提供する**: 追加前提（プロダクト側では API キーを提供しない）と矛盾し、運営コストも発生する
- **ローカル LLM 必須化**: non-goals.md に明示されているとおり採用しない

## 関連ドキュメント

- `docs/product/non-goals.md`（local LLM 必須化はしない）
- `docs/implementation/auth/auth-flow.md`（ユーザー自身の OpenAI アカウント利用を前提とした認証設計）
- `docs/implementation/ai/extraction-specification.md`（誤記修正が必要。Anthropic SDK → OpenAI SDK）
