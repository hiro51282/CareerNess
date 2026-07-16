#!/usr/bin/env bash
# パッケージングに同梱するリソース（Go server / web dist / codex）を resources/careerness に集める。
# electron-builder の extraResources がこのディレクトリを resourcesPath/careerness へ配置する。
set -euo pipefail

DESKTOP_DIR="$(cd "$(dirname "$0")/.." && pwd)"   # apps/desktop
IMPL_DIR="$(cd "$DESKTOP_DIR/../.." && pwd)"      # implementation
RES="$DESKTOP_DIR/resources/careerness"

rm -rf "$DESKTOP_DIR/resources"
mkdir -p "$RES"

echo "==> Go server をビルド"
EXT=""
if [ "$(go env GOOS)" = "windows" ]; then EXT=".exe"; fi
(cd "$IMPL_DIR/apps/api" && go build -o "$RES/server$EXT" ./cmd/server)

echo "==> web dist をコピー（要: 事前の pnpm build:web）"
if [ ! -f "$IMPL_DIR/apps/web/dist/index.html" ]; then
  echo "ERROR: apps/web/dist がありません。先に 'pnpm build:web' を実行してください" >&2
  exit 1
fi
cp -r "$IMPL_DIR/apps/web/dist" "$RES/web"

echo "==> codex CLI を同梱（CODEX_SRC 上書き可・無ければ同梱なしで続行）"
# npm インストールの codex は JS ラッパー（実体は optional dep の vendor 配下の
# native バイナリ）。ラッパーを同梱しても他マシンで動かないため、native を解決する。
# 本番配布では GitHub Releases のプラットフォーム別バイナリ取得に置き換える想定。
resolve_codex_native() {
  local real
  real="$(readlink -f "$1")"
  if ! head -c 2 "$real" | grep -q '#!'; then
    echo "$real" # 既に native バイナリ
    return 0
  fi
  local pkgdir
  pkgdir="$(cd "$(dirname "$real")/.." && pwd)" # @openai/codex パッケージルート
  find "$pkgdir" -type f \( -name codex -o -name codex.exe \) -path "*vendor*" 2>/dev/null | head -1
}

CODEX_SRC="${CODEX_SRC:-$(command -v codex || true)}"
NATIVE=""
if [ -n "$CODEX_SRC" ] && [ -e "$CODEX_SRC" ]; then
  NATIVE="$(resolve_codex_native "$CODEX_SRC")"
fi
if [ -n "$NATIVE" ] && [ -f "$NATIVE" ]; then
  cp "$NATIVE" "$RES/codex"
  chmod +x "$RES/codex"
  cp "$DESKTOP_DIR/third_party/codex/LICENSE" "$RES/codex-LICENSE"
  echo "    bundled(native): $NATIVE"
else
  echo "    codex の native バイナリが見つからないため同梱をスキップ"
fi

echo "==> staging 完了: $RES"
ls -la "$RES"
