<!--
CareerNESS PR テンプレート。Loop Engineering の人手 Checker ゲート。
不要な節は削除して構いません。チェックリストは既存の設計不変則（AGENTS.md）由来です。
-->

## 概要

<!-- 何を・なぜ変えたか。1〜3行で。 -->

**種別**: <!-- feat / fix / docs / ci / refactor / test のいずれか -->

## マージ権限の分類

この PR は次のいずれかに該当しますか？（該当する場合は**ユーザー承認**が必要。非該当なら Claude 自己マージ可）

- [ ] 重大な変更（広範囲・不可逆・高影響）
- [ ] セキュリティ要因（workspace 境界・認証・権限・API キー/鍵・公開範囲・CodeQL 指摘）
- [ ] 仕様変更に関わる判断（設計・ADR・データ契約・プロダクト方針）
- [ ] その他ユーザーの判断/承認が必要と考えられる事項

→ **いずれも未チェックなら「ルーティン PR」**（CI 緑で自己マージ対象）。

## DoD（完了条件）

- [ ] Go: `go build ./... && go vet ./... && go test ./...` が緑（`implementation/apps/api`）
- [ ] Web: `pnpm --filter web run test` と `pnpm build:web`（`tsc` 込み）が緑（`implementation`）
- [ ] CI（api / web / CodeQL）が緑

## Local-first / 設計不変則チェック

<!-- AGENTS.md の不変則。該当変更が無い場合もレビューで確認する。 -->

- [ ] CareerVault の facts / profiles / exports を DB に書いていない
- [ ] workspace 外のパスを参照・patch 対象にしていない
- [ ] ユーザー未承認の自動 apply が無い（変更は patch proposal → 承認 → apply を経由）
- [ ] session 内の inferred state を hidden truth 化していない
- [ ] export 文面を facts に逆流させていない
- [ ] `packages/*` が `apps/*` を import していない

## 関連

<!-- 関連 ADR / docs / Issue / 前後の PR など -->
