# 008 Desktop First — Desktop Host に Electron を採用

## Status

Accepted（PoC ゲート付き。PoC で致命的問題が出た場合は Tauri v2 へ再評価）

## Context

CareerNESS Editor のプロダクト原則の積み上げ——

- Career Vault はローカル資産（ADR-003）
- CareerNESS は AI credential を管理しない（ADR-007）
- AI はユーザー自身の認証済み環境（codex CLI）を利用する
- Structured Extraction を正式経路とする
- 将来的に Workspace Agent も視野に入れる

——の帰結として、**Desktop First** を現時点のプロダクト方針とした（`ai-foundation-direction.md`）。
ブラウザ単体ではユーザー認証済みの CLI をプロセスとして起動できず、現行実装の Vault 書き込みも
File System Access API（Chromium 限定）に依存している。

Desktop Host の候補は Electron / Tauri v2 / Wails v2 / PWA 継続。選定の評価軸は
**長期生存性・auto-updater・描画一貫性 ＞ フットプリント**とした。「重さ」は CareerNESS の
使われ方（**セッション型のエディタ**。常駐しない・ウィンドウ 1 枚・個人利用）に即して評価した。

## Decision

**Desktop Host として Electron を採用する。**

### アーキテクチャ原則: main は薄く保つ（thin-main）

`packages/shared` の「ゴミ箱にしない」と同型の不変則として定める:

> **main への追加は「Electron でしかできないこと」に限る。ロジックは Go に置く。**

- **main に置く**: ウィンドウ生成/ライフサイクル、Go サーバの子プロセス管理（ポート選択・
  ヘルスチェック・停止）、auto-updater（electron-updater）、ネイティブダイアログの IPC、
  セキュリティ設定（contextIsolation 等）
- **main に置かない**: ビジネスロジックと career data に触れる処理すべて
  （抽出 / patch / apply / `ResolveWithin` / codex exec / session / AI 状態診断は既存の Go 資産のまま）

### 構成

- BrowserWindow は **Go サーバが配信する `http://127.0.0.1:<port>`**（SPA 静的配信 + `/api/v1`）を
  読み込む。フロントエンド・API 契約は**無変更で再利用**（同一オリジン・CORS 変更なし）
- Vault 書き込みは Task2 で実装済みの Go 経路（session 束縛 + `ResolveWithin`）を正規経路へ昇格する
  （Chromium 同梱により File System Access API も動作するため、移行は段階的に行える）

### PoC ゲート

最小 PoC（ウィンドウ + Go 子プロセス起動 + attach → チャット → apply + codex exec）で検証し、
致命的問題が出た場合は **Tauri v2（Go sidecar）** へ再評価する。

## 重さの許容（判断根拠）

| 次元 | 実数 | セッション型エディタでの実害 |
|---|---|---|
| ダウンロード | インストーラ 80–150MB | なし（1 回きり。Slack / VS Code と同格） |
| ディスク | 展開後 ~250–350MB | なし |
| RAM | 実使用 ~300–400MB | 許容（使用中のみ。常駐アプリではない） |

加えて、ソロ開発では **3 系統の system webview（WebKitGTK / WebView2 / WKWebView）の互換テスト**を
保守するより、**単一 Chromium の更新追従**（electron-updater + CI で自動化が確立）の方が
総保守コストが低いと評価した。なお Electron ランタイム自体は重いままであり、
薄いのは**アプリとしてのアーキテクチャ層**である（この区別を明示的に許容した）。

## Consequences

### Positive

- 全 OS 同一描画（テストは 1 エンジンのみ）
- electron-updater による配布・自動更新の確立（Chromium セキュリティ更新の出荷義務を果たす道具が揃う）
- File System Access API が動作するため、Go-apply への移行を段階的に進められる
- アプリ内 `codex login` 起動（`codex-cli-integration.md` 課題#1 の残件）と
  codex CLI 同梱の検討（課題#2）の土台になる

### Negative

- 配布サイズ・メモリは Tauri / Wails 比で大きい（上記のとおり許容）
- Chromium セキュリティ更新に伴う定期リリースの義務を負う
- Node main 層が新設される（thin-main 原則で最小に保つ）

## Rejected Alternatives

- **Tauri v2**: 軽量・updater 内蔵・Go sidecar は公式パターンだが、Rust がスタックに加わり、
  3 系統の webview 互換（特に Linux の WebKitGTK）をソロで保守する負担が大きい。
  **PoC 破綻時の fallback として残す**
- **Wails v2**: Go 一本化・既存 mux の直マウントは魅力だが、auto-updater が DIY、
  エコシステムが小さく、v3 の alpha が長期化しているリスクを重く見た
- **PWA / ブラウザ継続**: ユーザー認証済み CLI を起動できず、credential 非管理原則（ADR-007）と
  両立しない。Vault 書き込みも Chromium 限定のまま

## 関連ドキュメント

- `docs/implementation/decisions/007-ai-credential-non-management.md`（Desktop First の駆動要因）
- `docs/implementation/decisions/001-workspace-boundary.md` / `003-local-first-storage.md` /
  `006-workspace-isolation-enforcement.md`（Go 経路昇格の基盤）
- `docs/implementation/ai/ai-foundation-direction.md`（Desktop First の初出・設計メモ）
- `docs/implementation/ai/codex-cli-integration.md`（課題#1 残件・課題#2 配布）
