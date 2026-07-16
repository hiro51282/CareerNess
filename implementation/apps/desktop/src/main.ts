// CareerNESS Desktop shell — thin-main（ADR-008）。
// ここに置くのは「Electron でしかできないこと」のみ:
// ウィンドウ / Go サーバ子プロセス管理 / ネイティブダイアログ IPC / セキュリティ設定。
// ビジネスロジック（抽出・patch・apply・codex 実行・診断）はすべて Go 側にある。
import { app, BrowserWindow, dialog, ipcMain } from 'electron'
import { spawn, type ChildProcess } from 'node:child_process'
import net from 'node:net'
import path from 'node:path'
import fs from 'node:fs'

let serverProc: ChildProcess | null = null
let mainWindow: BrowserWindow | null = null
let shuttingDown = false

// 多重起動を防ぐ（Go サーバ二重起動・vault 同時書き込みの回避）。
if (!app.requestSingleInstanceLock()) {
  app.quit()
}

/** 空きポートを OS に割り当てさせて返す。 */
function findFreePort(): Promise<number> {
  return new Promise((resolve, reject) => {
    const srv = net.createServer()
    srv.once('error', reject)
    srv.listen(0, '127.0.0.1', () => {
      const address = srv.address() as net.AddressInfo
      srv.close(() => resolve(address.port))
    })
  })
}

/** Go サーババイナリと Web dist の場所を解決する（packaged / dev で切替）。 */
function resolvePaths(): { serverBin: string; webDist: string } {
  if (app.isPackaged) {
    const res = path.join(process.resourcesPath, 'careerness')
    return { serverBin: path.join(res, 'server'), webDist: path.join(res, 'web') }
  }
  // dev: implementation/ からの相対（env で上書き可）
  const repo = path.join(__dirname, '..', '..', '..')
  return {
    serverBin: process.env.CAREERNESS_SERVER_BIN ?? path.join(repo, 'apps', 'api', 'bin', 'server'),
    webDist: process.env.CAREERNESS_WEB_DIST ?? path.join(repo, 'apps', 'web', 'dist'),
  }
}

/** /health が 200 を返すまでポーリングする。 */
async function waitForHealth(port: number, timeoutMs = 15000): Promise<void> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    try {
      const res = await fetch(`http://127.0.0.1:${port}/health`)
      if (res.ok) return
    } catch {
      // まだ起動中
    }
    await new Promise((r) => setTimeout(r, 150))
  }
  throw new Error('Go サーバが起動しませんでした（health check タイムアウト）')
}

/** Go サーバを子プロセスとして起動し、health 確認済みのポートを返す。 */
async function startServer(): Promise<number> {
  const { serverBin, webDist } = resolvePaths()
  if (!fs.existsSync(serverBin)) {
    throw new Error(`Go サーババイナリが見つかりません: ${serverBin}`)
  }
  if (!fs.existsSync(path.join(webDist, 'index.html'))) {
    throw new Error(`Web dist が見つかりません（先に build:web を実行）: ${webDist}`)
  }

  const port = await findFreePort()
  serverProc = spawn(serverBin, [], {
    // EXTRACTION_PROVIDER / CODEX_CLI_* 等はユーザー環境から透過する。
    env: { ...process.env, PORT: String(port), CAREERNESS_WEB_DIST: webDist },
    stdio: ['ignore', 'pipe', 'pipe'],
  })
  serverProc.stdout?.on('data', (d: Buffer) => console.log(`[go] ${String(d).trimEnd()}`))
  serverProc.stderr?.on('data', (d: Buffer) => console.error(`[go] ${String(d).trimEnd()}`))
  serverProc.on('exit', (code) => {
    serverProc = null
    if (!shuttingDown) {
      dialog.showErrorBox('CareerNESS', `バックエンドが予期せず終了しました (exit=${code})`)
      app.quit()
    }
  })

  await waitForHealth(port)
  return port
}

/** Go サーバを確実に停止する（SIGTERM → 猶予後 SIGKILL）。 */
function stopServer(): void {
  shuttingDown = true
  const proc = serverProc
  if (!proc) return
  serverProc = null
  proc.kill('SIGTERM')
  const killTimer = setTimeout(() => {
    try {
      proc.kill('SIGKILL')
    } catch {
      // 既に終了済み
    }
  }, 2000)
  proc.once('exit', () => clearTimeout(killTimer))
}

async function createWindow(port: number): Promise<void> {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
    },
  })
  // フロントと API は Go サーバの同一オリジン（ADR-008）。
  await mainWindow.loadURL(`http://127.0.0.1:${port}`)
  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

// ネイティブのフォルダ選択（attach 用）。絶対パスを返す＝Go 経路 attach の入力になる。
ipcMain.handle('careerness:pickDirectory', async (): Promise<string | null> => {
  const result = await dialog.showOpenDialog(mainWindow ?? undefined!, {
    title: 'Career Vault フォルダを選択',
    properties: ['openDirectory', 'createDirectory'],
  })
  return result.canceled || result.filePaths.length === 0 ? null : result.filePaths[0]
})

app.whenReady().then(async () => {
  try {
    const port = await startServer()
    await createWindow(port)
  } catch (e) {
    dialog.showErrorBox('CareerNESS 起動エラー', String(e))
    app.quit()
  }
})

app.on('window-all-closed', () => {
  app.quit()
})

app.on('will-quit', () => {
  stopServer()
})

process.on('exit', () => {
  // 最後の保険（通常は will-quit で停止済み）。
  try {
    serverProc?.kill('SIGKILL')
  } catch {
    // noop
  }
})
