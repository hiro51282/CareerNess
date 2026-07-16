# Electron PoC レポート（ADR-008 ゲート判定）

実施: 2026-07-16〜17 / ブランチ: `poc/electron-shell` / 環境: Ubuntu (Wayland), Electron 43.1.1, electron-builder 26.15.3

## 判定: **GO**（Electron 採用を維持。致命的問題なし・Tauri 再評価は不要）

※ win/mac ビルド成立（ゲート条件4）のみ、workflow_dispatch の制約上 **master マージ後に実行**して確定する（後述）。

## ゲート条件と結果

| # | 条件 | 結果 |
|---|---|---|
| 1 | ウィンドウ＋Go 子プロセスで attach→チャット→Go 経路 apply→Facts 反映 | ✅ SPA 配信→attach（絶対パス）→mock 抽出→`/apply-patch`(200)→`/workspace/files` 反映を確認。ウィンドウ起動・SPA 読込も確認 |
| 2 | 終了時にプロセスリークなし | ✅ electron 本体へ SIGTERM → will-quit → Go 子プロセスも終了（dev/packaged 両方で確認） |
| 3 | Linux パッケージから同梱 codex で実 AI | ✅ パッケージ版で `ai/status` が `ready:true "Logged in using ChatGPT"`（同梱 native バイナリ・絶対パス注入・`~/.codex` 認証共有） |
| 4 | win/mac 未署名ビルド成立 | ⏳ workflow（manual dispatch）追加済み。**マージ後に実行して確定**（実機起動検証は配布タスク） |
| 5 | 致命的問題ゼロ | ✅ なし |

## 実測値

| 項目 | 値 |
|---|---|
| AppImage（codex 同梱） | **208MB** |
| 展開後 | **530MB**（うち codex 209MB） |
| AppImage（codex 非同梱の参考値） | 127MB / 展開 321MB |
| Go server バイナリ | 9.9MB |
| 起動（electron→Go health→SPA 表示） | 体感 2〜3 秒 |

ADR-008 の重さ予測（インストーラ 80–150MB／展開 250–350MB）は **codex 非同梱なら的中**。
**codex 同梱で約 +200MB** が新事実（下記 発見-1）。

## P0 検証結果

- **P0-2（GUI 起動の PATH 問題）**: **解決**。同梱 codex を `CODEX_CLI_BIN` に**絶対パス注入**するため PATH 非依存。
  ※ 同梱しない構成で system codex に頼る場合は依然 PATH 問題が残る（配布タスクの考慮点）。
- **P0-3（子プロセスライフサイクル）**: **成立**。SIGTERM→graceful stop（SIGTERM→2s→SIGKILL）。
  注意: 開発時に **pnpm ラッパーを kill すると electron に伝播せず残存**する（electron 本体を signal すること）。
- **P0-4（codex の Windows 対応）**: **対応あり**。npm ラッパーの対応表に `@openai/codex-win32-x64 / win32-arm64`
  （`x86_64-pc-windows-msvc` 等）が存在 → **Windows ネイティブバイナリは配布されている**。同梱可能。
- **P0-1（Electron での File System Access API）**: **検証不要化**。desktop モードを最初から Go 経路
  （attach=ネイティブダイアログ→session 束縛 / read=`/workspace/files` / apply=`/apply-patch`）で実装したため、
  ブラウザ FS への依存自体が無い。ADR-008 の「段階的移行」仮定は使わずに済んだ。

## 発見事項（配布タスクへの引き継ぎ）

1. **codex native バイナリは 209MB**。同梱するとパッケージがほぼ倍増（127→208MB）。
   配布タスクで「(a) 同梱（インストール不要 UX 優先）vs (b) 初回起動時ダウンロード」を実数ベースで判断する。
2. **npm の `codex` コマンドは JS ラッパー**（シンボリックリンク→codex.js）。実体は optional dep
   `@openai/codex-<platform>` の vendor 配下。**ラッパーだけ同梱しても他マシンでは動かない**
   （stage.sh は native 解決を実装済み。本番配布は GitHub Releases / npm からのプラットフォーム別取得に置き換える）。
3. **localhost API 保護**は動的ポートのみ実装。**共有トークンは配布前必須**（P2 として据え置き）。
4. thin-main 原則は維持できた（main.ts ~190 行。ロジック追加はゼロ、env 注入とプロセス管理のみ）。

## 残作業（PoC 外）

- win/mac ビルドの実行（マージ後 `gh workflow run desktop-build.yml`）と Windows 実機起動検証（実機入手後）
- 配布タスク: codex 同梱方式の決定・署名・auto-updater・アプリ内 `codex login` 起動・localhost トークン
- デスクトップ UI 磨き: WorkspaceAttach の文言（「Chrome/Edge が必要」は desktop で不適切）等

## 構成（記録）

```
Electron main（thin, apps/desktop）
  └─ 空きポート探索 → Go server 子プロセス（PORT/CAREERNESS_WEB_DIST/CODEX_CLI_BIN 注入）
       └─ SPA 配信（/) ＋ /api/v1（同一オリジン・127.0.0.1 bind）
  └─ BrowserWindow → http://127.0.0.1:<port>
  └─ preload: desktop フラグ + pickDirectory IPC のみ
フロント desktop モード: attach=ダイアログ→/workspace/attach、read=/workspace/files、apply=/apply-patch
```
