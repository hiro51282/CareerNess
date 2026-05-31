# Reviewer向け引継ぎメモ

作成者: PM エージェント (Claude Sonnet 4.6)  
作成日: 2026-05-31  
ブランチ: agent-planner  
レビュー対象: `docs/agent-plan.md`

---

## このメモの目的

`docs/agent-plan.md`（PM分析・ロードマップ）をReviewer Agentへ引き継ぐ。  
計画の判断根拠・不確実な箇所・懸念点を伝え、独立したレビューを促す。

---

## この計画で重要視した点

### 1. Local-first 原則との整合性を各フェーズで維持すること

マルチユーザー対応（Feature A）でプラットフォーム DB（Supabase）を導入するが、**career data（facts/profiles/exports）は DB に入れない** という分離が崩れやすいポイント。`docs/architecture/platform-vs-vault.md` が定義する「DB = platform metadata、CareerVault = career truth」の境界を、ロードマップ各フェーズのアーキテクチャ変更ポイント（section 8）で明示した。

### 2. フェーズ順序に依存関係の根拠を持たせること

A → B → C の順序を直感ではなく理由付きで提示した。

- **A（マルチユーザー）が先**：認証なしで AIキャリア分析を公開すると API が無保護になり、OpenAI コストが制御できない
- **B（AIキャリア分析）が先**：面接支援の回答品質は accumulated facts の量と質に依存する。facts が少ない状態で面接支援を先に作っても品質が低い

### 3. AI面接支援もパッチ提案モデルに従わせること

面接回答案・準備シートを「直接保存するUI」にすると ADR-002（Patch-Oriented Editing）が崩れる。section 8 で明示したとおり、面接支援の output は必ず `exports/interview-prep/` への Patch proposal として提示し、ユーザー承認後に保存する設計とした。

### 4. 技術的負債を現時点の実装と照合して記録したこと

`IMPLEMENTATION_PROGRESS.md` のリストをそのまま転記せず、実際のソースコード（`internal/extraction/models.go` と `internal/patch/model.go`、`docs/implementation/workspace/fact-schema.md` と YAMLFact struct）を照合して新たな負債を発見・追記した。

---

## 不確実な前提

### U-1: Codex AppServer の仕様が不明瞭

`docs/ai/ai-behavior.md` と `CLAUDE.md` には「Codex AppServer 経由で AI を呼び出す」と書かれているが、**Codex App Server がどのようなプロトコルでリクエストを受け付けるかの仕様書が存在しない**。

一方 `docs/implementation/ai/extraction-specification.md` には `github.com/anthropic-ai/sdk-go` を使って直接 Claude API を呼ぶ実装例が示されており、「Codex AppServer 経由」と矛盾している。

この乖離が CodexExtractionProvider 実装時に何を呼ぶのかを曖昧にしている。ロードマップでは「Codex AppServer 経由」を前提としたが、実際は直接 Claude API 呼び出しに変更になる可能性がある。

### U-2: Phase 1 の工数は楽観的かもしれない

「AI 実統合 L（1週間）」と見積もったが、上記 U-1 の Codex AppServer 問題が解決されていない場合、実装方針の決定だけで時間を使う可能性がある。

### U-3: マルチデバイス sync は設計が存在しない

Feature A-3（マルチデバイス対応）で「Git-backed sync を選択肢として提案」としたが、これは CareerVault を Git リポジトリとして管理するという提案であり、**実際の UX 設計が存在しない**。ユーザーが Git に慣れているという前提が不明確。

### U-4: Feature B/C の AI 品質は未知数

キャリア分析の「強み・パターン抽出」や面接支援の「STAR 回答生成」の品質は、プロンプト設計と facts のデータ質に強く依存する。機能として完成していても「精度が低くて使えない」状態になりうる。プロンプトエンジニアリングのコストが計画に含まれていない。

### U-5: File System Access API の長期依存

現在のブラウザ側 workspace attach は Chrome/Edge のみ対応。**将来的に Electron アプリ化するのか、CLI ヘルパー経由にするのか、ブラウザ API のまま行くのかが未決定**。この決定がなければ「Firefox/Safari でも動く」という要求が出た際に設計の見直しが発生する。

---

## レビューで重点的に見てほしい箇所

### R-1: section 8 のアーキテクチャ変更ポイント

マルチユーザー対応時に「Supabase には career data を入れない」という制約が、提示した DB スキーマ案（users, sessions, workspace_attachments, usage_logs テーブル）で本当に守れるかを確認してほしい。特に `workspace_attachments` テーブルに workspace 内のファイルパス以上の情報を持たせてしまうケースが起きやすい。

### R-2: section 5（技術的負債）の優先度判断

「Patch struct の二重定義」と「API 認証なし」を高優先度（MVP 前に対処すべき）に分類した。この分類が Phase 1 の残りタスクとして実装チームにとって現実的かを確認してほしい。Patch struct 統合は既存テストの修正が伴う可能性がある。

### R-3: Feature A と B の並行着手可否

ロードマップ note に「Feature A と Feature B の一部は並行着手可能」と書いたが、**AI 実統合（CodexExtractionProvider）と認証基盤（Supabase Auth）を同時に進める場合、同一 Go API サーバーへの変更が競合しやすい**。並行着手を推奨するなら feature branch 戦略を明示すべきかもしれない。

### R-4: AI面接支援の Patch proposal 適用が UX として成立するか

面接回答案を patch proposal として提示する設計は、設計思想としては正しいが、**「模擬面接中に都度 patch を承認するフロー」が UX として受け入れられるか**が不明。面接モードは会話テンポが重要で、patch review の割り込みが体験を壊す可能性がある。Phased-approval（面接後にまとめて確認）などの変形が必要かもしれない。

### R-5: section 9（優先度判断の指針）の Local-first チェックリスト

このチェックリストが実装レビュー時に実際に使われる形式になっているかを確認してほしい。現在は箇条書きだが、PR テンプレートや CI チェックとして組み込む形にすべきか検討の余地がある。

---

## 自分で懸念している点

### C-1: Codex AppServer vs. 直接 Claude API 呼び出しの乖離が最大のリスク

`extraction-specification.md` に書かれた Claude API 直接呼び出しの実装例と、CLAUDE.md の「Codex AppServer 経由」が矛盾している。この計画では「Codex AppServer 経由」を前提にロードマップを書いたが、もし実装チームが直接 Claude API 呼び出しに舵を切ると、**ユーザーが自分の OpenAI アカウントを使うという前提（非 Anthropic API）と矛盾する**。この前提の揺らぎが Phase 1 の AI 実統合タスクに影響する。

### C-2: packages/* の分離タイミングが遅すぎるかもしれない

現在すべての実装が `apps/api` 内に集中しており、`packages/workspace-core`・`packages/patch-engine`・`packages/schema` は README のみ。Phase 2 以降で分離するとしたが、**Feature A（マルチユーザー）・Feature B（AIキャリア分析）と並行して実装が積み上がると、後から packages に分離するコストが高くなる**。Phase 1 完成後・Phase 2 着手前に「packages 分離スプリント」を入れるべきかもしれない。

### C-3: 面接支援フェーズが2026-12末まで後ろ倒しになっていること

競合サービスが同様機能をリリースするリスクがある。もし「AI面接支援」がプロダクトの差別化要素として最重要なら、Phase 順序を見直し（facts の minimum viable 状態でも面接支援に入れるか検討）か、Feature B と Feature C を並行開発できるよう分割できる機能がないかを再検討すべきかもしれない。

### C-4: 認証設計が「方針のみ」で具体実装が書かれていない

`docs/implementation/auth/auth-flow.md` は「やること・やらないこと」は書いているが、JWT の具体的なクレーム設計、Supabase Auth との接続方法、Go 側のミドルウェア実装方針が存在しない。Feature A 着手時に実装者が設計から始めることになり、想定工数より時間がかかる可能性がある。

### C-5: Fact schema の不整合は「Phase 2 で解消」が正しいか

docs の `fact-schema.md`（`action/decision/impact/context/evidence` フィールドあり）と実際の `YAMLFact` struct（`description` に集約）の乖離を「Phase 2 で解消」とした。しかし、**Phase 1 でユーザーが facts を蓄積し始めた後に schema を変更すると、既存データのマイグレーションが必要になる**。Phase 1 完了前に最低限の schema を fix すべきか、再検討を促したい。

---

## 引継ぎ先への依頼

上記の懸念点を踏まえ、以下の観点でレビューをお願いしたい。

1. **U-1（Codex AppServer 仕様の乖離）** に対して、実際のコードベースと他ドキュメントから真の意図を読み取り、方針を明確化してほしい
2. **C-2（packages 分離のタイミング）** について、Phase 1 完成前・後どちらが適切かの意見を
3. **R-4（面接支援の UX）** について、patch proposal model を守りつつ会話テンポを損なわない代替設計案があれば提案してほしい
4. **C-5（fact schema 不整合）** の対処タイミングについて、データマイグレーションコストを考慮した判断を
