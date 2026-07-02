# AI 基盤の方向性（設計メモ）

> 本書は ADR ではなく**設計意図の記録**。確定アーキテクチャ決定（却下案を伴う ADR）ではない。
> ただし本メモの「プロダクト原則」「Desktop First」は方針が固まったため、近く ADR 化する候補
> （ADR-004「AI Provider Integration」の改訂/supersede を含む。§末尾参照）。
> 初版記録時点: Task4 Phase B 着手前 / 直近更新: プロダクト方針転換（Credential 非管理・Desktop First・Codex CLI 正式経路）を反映。

CareerNESS はまず **Career Vault を編集・管理するための Editor（CareerNESS Editor）** を完成させることを
MVP の目的とする。Export / Resume / 将来の Jobs / AI 分析・RAG は、その上に載る利用形態である。
本メモはその位置づけを前提に、AI 基盤の方針を固定する。

---

## プロダクト原則：CareerNESS は AI Provider の Credential を管理しない

**これは MVP の都合ではなく、CareerNESS のプロダクト原則である。** 技術的理由ではなく
**信頼境界（trust boundary）を明確にするため**の判断。

- CareerNESS は **API キーを保持しない**。
- **API キー入力を基本 UX にしない**。
- **AI 利用契約はユーザー自身が管理する**（ユーザーが自分で認証済みの AI 実行環境を使う）。

将来的に API Provider 系統を完全排除すると決めるものではない（§「Codex CLI を正式経路とする」参照）。
MVP では「**Credential を CareerNESS が管理しない**」ことを設計原則として固定する。

## Desktop First（現時点のプロダクト方針）

**Desktop First は UI 技術（Electron / Tauri / Wails 等）の決定ではなく、CareerNESS Editor の
プロダクト方針である。** Electron ありきで決めたのではなく、次のプロダクト原則を積み上げた
帰結として位置づける：

- Career Vault はローカル資産である
- CareerNESS は AI Provider の Credential を管理しない
- AI はユーザー自身の認証済み環境（CLI 等）を利用する
- Structured Extraction を正式経路とする
- 将来的に Workspace Agent も視野に入れる

これらを満たすと、現時点では **Desktop ＋ ローカル CLI 利用** が最も自然な構成であり、
**現時点では他に同等の要件を満たす現実的な選択肢が見当たらない**（例：ブラウザ単体では、
ユーザーが認証済みの CLI をプロセスとして起動できない）。ゆえに Desktop First を
**現時点のプロダクト方針**として採用する。

一方で、**Desktop Host の実装技術（Electron / Tauri / Wails 等）は現時点では決定しない**。
ここは別判断として残す。移行は「作り直し」ではなく、**既存の React / Go / Provider 資産を
最大限活かす**前提で行う。

---

## Structured Extraction：正式経路は Codex CLI

CareerNESS の AI 抽出は、当面 **Structured Extraction Provider** 方式のみを正式経路とする。
LLM を「structured JSON extraction provider」として扱い、会話 → 構造化 JSON → Go 側で
deterministic に validate / normalize / schema 検証する既存パイプライン（`extraction-specification.md`、
ADR-002 のパッチ提案モデル）と整合させる。

`ExtractionProvider` interface（`ExtractFacts(ctx, conversation) → ExtractedFactResult` と `Name()`）を
中心に据える。MVP の provider 構成：

- **Mock**（既定）：開発・CI が AI 実行環境なしで会話 UX を確認できるようにする。
- **Codex CLI（正式な実 AI 経路）**：ユーザーが独立に認証済みの `codex` CLI を実行し、
  構造化 JSON を得る。**Credential は CLI 側に閉じ、CareerNESS は鍵に触れない**（上記プロダクト原則と整合）。

### HTTP（OpenAI API 直叩き）Provider は現行 MVP の正式経路から外す（deprecated）

現行の HTTP `CodexExtractionProvider`（OpenAI API を `OPENAI_API_KEY` で直接呼ぶ実装）は、
**CareerNESS が Credential を管理する前提**で設計されているため、**現行 MVP の正式経路から外す**（deprecated 扱い）。

- これは「**HTTP という通信方式を否定するものではない**」。**Credential 管理方針との整合による優先順位付け**である。
- 実装は削除せず休眠保持する。将来 CareerNESS が十分成熟し API 提供を正式サポートする段階になれば、
  HTTP Provider 系統を**復活・再設計する可能性を残す**。

## Runtime 抽象化は行わない（YAGNI）

現時点では Claude Code / Gemini CLI 等まで見据えた抽象化は**行わない**。**Codex CLI のみ**を実装する。

- **Runtime Interface / Executor Interface は新設しない。**
- ただし、`CodexCLIProvider` の内部で **CLI 実行処理だけは Prompt 生成・Extraction 固有処理と疎結合**に保つ。
  具体的には「CLI 実行」を、`stdin(string) → stdout(string)` のみを知り、prompt 文言・`ExtractedFactResult`・
  JSON schema を**知らない**小さな内部関数（例：`runCLI`）に隔離する。
- **TODO（将来）**：CLI 実行処理は、Claude Code / Gemini CLI 等を追加したくなった時点で
  **Runtime / Executor 層として独立させる可能性がある**。その際は `runCLI` 相当を昇格させれば足りる
  構造にしておく（今は interface を作らない）。

---

## 将来検討事項：Workspace Agent（MVP 対象外）

将来的には **Workspace Agent**（Codex CLI / Claude Code / Gemini CLI 等の、ツール実行を伴う
エージェント方式）も採用候補とする。ただし **MVP の対象外**とする。

> **Workspace Agent は、Structured Extraction Provider の置き換えではなく、別責務として共存することを想定する。**
> 既存の structured extraction 経路を廃止・代替するものではない。

理由：
- エージェント方式は、ツール実行・反復的な探索・ワークスペース横断の操作など、
  structured extraction とは**責務の大きさ・性質が異なる**。
- MVP のスコープ（会話 → 構造化抽出 → パッチ提案）を超える。
- 現時点で実装すること、または将来のために**過剰な抽象化を先取りすること**は YAGNI であり、
  本プロジェクトの方針（必要になった時点で導入）に反する。

## 将来追加する場合の前提：Provider の責務は変えない

Workspace Agent を追加する場合も、**既存 `ExtractionProvider` の責務は広げない**。

- `ExtractionProvider` は「会話 → `ExtractedFactResult`」の structured extraction に限定し続ける。
- Workspace Agent は Provider の拡張（interface への機能追加）ではなく、
  **orchestration 層における別の責務（新しい境界）として導入**する想定。
- これにより Provider interface を汚さずに両方式を併存させられ、structured extraction の
  deterministic なパイプライン（validate/normalize/schema）も影響を受けない。

## 求めるもの（チェックリスト）

- ✅ CareerNESS は AI Provider の Credential を管理しない（プロダクト原則）
- ✅ Desktop First は現時点のプロダクト方針（原則の積み上げの帰結）。Desktop Host の実装技術（Electron / Tauri / Wails 等）は未決定・別判断
- ✅ 正式経路は Structured Extraction → Codex CLI（HTTP は deprecated・休眠保持）
- ✅ Runtime 抽象化はしない（Codex CLI のみ／CLI 実行部分だけ疎結合を維持）
- ✅ Provider の責務を広げない（Workspace Agent は別責務として将来共存）

## ADR 化の予定（メモ → 確定判断）

以下は方針が固まったため、近く ADR 化を検討する（本メモは先行する設計意図の記録）：

- **Credential 非管理のプロダクト原則**：ADR-004「AI Provider Integration」（ユーザーの OpenAI 鍵を
  AppServer 経由で使う前提）を**改訂/supersede** する新 ADR。以前保留した「credential handling（ADR-007 相当）」は、
  この原則により方向が確定する（案C→案B の鍵リレー路線は破棄）。
- **Desktop First**：Desktop Host 採用の背景と、具体実装（Electron/Tauri）を別判断とする旨を記録する ADR。
