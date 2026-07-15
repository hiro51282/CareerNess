# 007 CareerNESS は AI Provider の credential を管理しない

## Status

Accepted（ADR-004 の credential 取り扱い部分を改訂する）

## Context

ADR-004（2026-05-31）は「AI 呼び出しはユーザー自身の OpenAI アカウント（Codex 経由）を使い、
CareerNESS 運営は請求主体にならない」と決定したが、**ユーザーの API key を AppServer が
どう受け取るか**（一時受け取り方式 vs ブラウザ直接呼び出し）は未確定のまま残されていた。

Task4 Phase A では暫定的に server env key 方式（`OPENAI_API_KEY` を env から読む HTTP provider）を
実装したが、これは CareerNESS が credential を扱う構造であり、ローカル単一ユーザーの暫定措置
だった。

その後のプロダクト方針の再整理（2026-07、`ai-foundation-direction.md`）で、これは技術選択では
なく**信頼境界（trust boundary）の問題**であると整理された。CareerVault（キャリアの正本）を
預からないのと同じ理由で、AI の credential も預からないことがプロダクトの信頼性の核になる。

実証として、credential-free の `CodexCLIProvider`（ユーザーが独立に認証済みの codex CLI を
非対話実行）が実装され（PR #19）、実 AI で end-to-end 動作を確認済み（2026-07-06 smoke、
`codex-cli-integration.md`）。

## Decision

**CareerNESS は AI provider の credential を保持・保存・中継しない。** これをプロダクト原則とする。

- CareerNESS（AppServer / UI / DB / ログ）は API キー・トークンに一切触れない
- **API キー入力を基本 UX にしない**
- AI 利用契約・認証は**ユーザー自身が管理**する（`codex login` 等。認証状態は CLI 側
  `~/.codex` に閉じる）
- 実 AI の正式経路は **Structured Extraction → Codex CLI**（`EXTRACTION_PROVIDER=codex-cli`）
- HTTP/OpenAI 直叩き provider（`OPENAI_API_KEY` 前提）は **deprecated・休眠保持**。
  これは HTTP という通信方式の否定ではなく credential 方針との整合による判断であり、
  将来 CareerNESS が API 提供を正式サポートする段階になれば再設計の余地を残す

### ADR-004 との関係

- **維持**: ユーザー自身のアカウントを使う／CareerNESS 運営は AI 利用料の請求主体にならない
- **改訂**: 「AppServer がブラウザとユーザーの OpenAI アカウントの中継層として機能する」前提と、
  Open Question「API key をどう受け取るか」は本 ADR が解決・置換する（受け取らない、が答え）

## Consequences

### Positive

- 信頼境界が明確になる（CareerNESS は正本もクレデンシャルも預からない）
- 鍵の漏洩・誤ログ・誤永続化のリスクがアプリから構造的に消える
- 「credential を保持するコード」を禁止パターンとして機械・レビューで強制できる
  （AGENTS.md / PR テンプレートに反映済み）

### Negative

- 実 AI の利用に codex CLI のインストールと `codex login` が前提になる
  （検知＋案内は実装済み: `GET /api/v1/ai/status`、PR #25。アプリ内から login を起動する動線は
  ADR-008 の Desktop Host 上で対応予定）
- ブラウザ単体では実 AI を利用できない（→ ADR-008 Desktop First の駆動要因）

## Rejected Alternatives

- **リクエスト毎の一時 key リレー（旧・案B）**: AppServer が一瞬でも鍵に触れる構造は信頼境界を
  曖昧にし、鍵入力 UX も必要になる
- **server env key の恒久化（旧・案C）**: 自己ホスト・単一ユーザーでしか成立せず、配布する
  プロダクトの前提にできない
- **ブラウザから AI API を直接呼ぶ（旧・案A）**: 信頼コア（validate/normalize）のクライアント
  二重化と、ブラウザへの鍵露出を招く

## 関連ドキュメント

- `docs/implementation/decisions/004-ai-provider-integration.md`（本 ADR が credential 部分を改訂）
- `docs/implementation/decisions/008-desktop-first-electron.md`（本原則の帰結としての Desktop First）
- `docs/implementation/ai/ai-foundation-direction.md`（原則の初出・設計メモ）
- `docs/implementation/ai/codex-cli-integration.md`（実証と運用手順）
