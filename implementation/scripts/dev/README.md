# Development Scripts

開発用コマンドスイート。リポジトリルートから実行してください。

## Scripts

### `./up.sh`

Go API サーバー（:8080）と React dev サーバー（:5173）を同時起動します。

```bash
./implementation/scripts/dev/up.sh
```

**出力例：**
```
🚀 Starting CareerNess development servers...

📡 Starting Go API server (:8080)...
⚛️  Starting React dev server (:5173)...

✅ Both servers running:
   • API  → http://localhost:8080
   • Web  → http://localhost:5173

Press Ctrl+C to stop all servers
```

**停止：** `Ctrl+C` で両サーバーが停止します。

---

## 将来のスクリプト候補

- `test.sh` — Go テスト + TypeScript 型チェック
- `build.sh` — React + Go の本番ビルド
- `lint.sh` — フォーマット・linting チェック
