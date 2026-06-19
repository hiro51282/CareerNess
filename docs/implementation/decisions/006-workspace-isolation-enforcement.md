# 006 Workspace Isolation Enforcement (MVP)

## Status

Accepted

## Context

ADR-001（Workspace Boundary）は「AI が参照・提案・patch 対象にできるデータは、ユーザーが明示的に attach した workspace 配下に限る」と*定めて*いるが、その境界を backend 側で*強制する機構*は MVP 時点で未実装だった。

現状コードには次の穴がある：

- `POST /api/v1/apply-patch`（`handler/apply.go`）は実際にサーバ FS へ書き込むが、`op.Target` を再検証しない。`patch.Applier.ApplyPatch` は `patch.Validate` を呼ばないため、`op.Target = "../../etc/passwd"` のような操作が `filepath.Join(workspacePath, op.Target)` で workspace 外へ書き込み得る。
- apply の root を、リクエストボディの任意 `workspace_path` から受け取って信用していた。allowlist が無く、プロセスが書ける任意の絶対パスを射程に入れてしまう。`filepath.Clean` + `Contains("..")` の検査は、絶対パスのエスケープ（例 `/srv/a/../../etc` → `/etc`）を取りこぼす。
- `session_id` は patch に echo されるだけの文字列で、session ↔ workspace_root の束縛が存在しない。どの session がどの root に書いてよいかをサーバが引けない。

なお File System Access API は絶対パスをブラウザに渡さない（`handle.name` は basename のみ）。現行 Chromium の**ライブ書き込みはブラウザ FS 経路**であり、その境界は OS が保証する `dirHandle` のスコープで守られている。`/apply-patch`（server-apply）は、サーバが vault と同一 FS を共有する構成（ローカル / dev / 将来のローカル helper）でのみ意味を持つ別経路である。

## Decision

ADR-001 の境界を backend 側で強制する最小機構を導入する。

- **server-apply（`/apply-patch`）は堅牢化して残す**。無効化はしない（無効化は Task2 の目的＝apply 経路の isolation 実現から外れる）。
- **session → workspace_root の束縛を最小フィールドで実装する**。Attachment が持つのは `session_id` / `workspace_id` / `workspace_root` の 3 つのみ。in-memory で保持する transient な写像であり、キャリアデータも workspace mirror も持たない。
- **apply は三段で root 外書き込みを不能化する**：
  1. **認可**: `session_id` 必須。対応する attachment が無ければ拒否。
  2. **構造検証**: `patch.Validate` を apply 入口で必ず通す。
  3. **封じ込め**: 各 operation の `target` を `workspace.ResolveWithin(root, target)` で root 配下に解決し、外れるものは拒否。
- **root は attachment（session store）から導出し、リクエストボディからは受け取らない**。`apply-patch` のボディは `{session_id, patch}` とし、`workspace_path` は廃止する。
- **境界判定は純関数 `workspace.ResolveWithin` に一元化する**。文字列パターン検査ではなく `filepath.Rel` ベースの封じ込め＋symlink 再評価で判定する。

## Consequences

### Positive

- backend の書き込みが attach された root 内に封じ込められ、ADR-001 が機構として強制される。
- 境界検証ロジックが単一ソース（`ResolveWithin`）になり、apply・将来の経路で再利用・テストできる。
- root をサーバ側 attachment から導出することで、ボディ経由の任意パス指定という攻撃面を排除できる。

### Negative

- `/apply-patch` の契約が変わる（`workspace_path` 廃止・`session_id` 必須）。ただし現状フロントは `/apply-patch` を呼んでおらず（ライブ書き込みはブラウザ FS 経路）、実害は小さい。
- server-apply はサーバが vault と同一 FS を共有する構成でのみ有効。クラウド AppServer 構成では使えない（その構成では backend は user vault への FS アクセスを持たない前提）。

## Rejected Alternatives

- **案B: `/apply-patch` を無効化する**: Task2 の目的（apply 経路の isolation）から外れる。非 Chromium ブラウザや将来のローカル helper のために server-apply を生かす余地も失う。
- **フル session モデル（user / workspace attachment / conversation / patch review の 4 種別）の先行実装**: MVP では Workspace Attachment しか必要なく YAGNI。`session-model.md` の理想形は将来実装として残す。
- **将来用フィールド（`UserID` / `ExpiresAt` / `Scope` / `Revision`）の先行追加**: 現時点で利用されない。認証・失効・stale 判定が必要になった Task3 以降で追加する。

## Task3 以降への引き継ぎ

- `UserID` による user スコープの attachment ルックアップ（JWT 認証導入時）。session_id だけで apply を許す現状は、認証導入までの単一ユーザー MVP 前提であることを明記する。
- attachment 失効（`ExpiresAt`）と stale 判定（`session-model.md` の staleness handling）。

## 関連ドキュメント

- `docs/implementation/decisions/001-workspace-boundary.md`（本 ADR が境界を*強制*する対象）
- `docs/implementation/decisions/003-local-first-storage.md`（backend は truth を持たない → Store は transient なポインタのみ）
- `docs/implementation/decisions/004-ai-provider-integration.md`（AI provider 統合とは独立した境界の話）
- `docs/implementation/backend/session-model.md`（MVP 実装スコープの注記を追記）
- `docs/implementation/backend/backend-structure.md`（apply は attachment 必須の注記を追記）
