import { createContext, useState, useCallback, type ReactNode } from 'react'
import type { Operation, Patch } from '../../types/patch'
import { attachWorkspace, getWorkspaceFiles, applyPatchViaServer } from '../../api/client'
import { upsertFactYaml, markFactStatusYaml } from './factFile'

// desktop（Electron）モード判定。preload が window.careerness を公開している時のみ true。
// desktop では attach/read/apply を Go 経路（ADR-006/008 の正規経路）で行い、
// ブラウザでは従来どおり File System Access API を使う。
const isDesktop = typeof window !== 'undefined' && window.careerness?.desktop === true

export interface WorkspaceContextType {
  dirHandle: FileSystemDirectoryHandle | null
  workspaceId: string
  files: Map<string, string>
  /** workspace が attach 済みか（browser: dirHandle / desktop: session）。 */
  attached: boolean
  attach: () => Promise<void>
  refreshFiles: () => Promise<void>
  /** 承認済み patch を workspace へ適用する（mode に応じ browser FS / Go 経路）。 */
  applyPatch: (patch: Patch) => Promise<string[]>
}

export const WorkspaceContext = createContext<WorkspaceContextType | null>(null)

export function WorkspaceProvider({ children }: { children: ReactNode }) {
  const [dirHandle, setDirHandle] = useState<FileSystemDirectoryHandle | null>(null)
  const [desktopSession, setDesktopSession] = useState<string | null>(null)
  const [workspaceId, setWorkspaceId] = useState('local-careervault')
  const [files, setFiles] = useState<Map<string, string>>(new Map())

  const attach = useCallback(async () => {
    if (isDesktop) {
      // desktop: ネイティブダイアログ → 絶対パスを Go に attach（session 束縛）→ Go 経由で読込
      const dir = await window.careerness!.pickDirectory()
      if (!dir) return
      const res = await attachWorkspace(dir)
      setDesktopSession(res.session_id)
      setWorkspaceId(res.workspace_id)
      const wf = await getWorkspaceFiles(res.session_id)
      setFiles(new Map(Object.entries(wf.files)))
      console.debug(`[workspace] attached via server: ${res.workspace_id}`)
      return
    }
    const handle = await window.showDirectoryPicker({ mode: 'readwrite' })
    console.debug(`[workspace] attached: ${handle.name}`)
    setDirHandle(handle)
    setWorkspaceId(handle.name)
    const loaded = await readWorkspaceFiles(handle)
    console.debug(`[workspace] files loaded: ${loaded.size}`, [...loaded.keys()])
    setFiles(loaded)
  }, [])

  const refreshFiles = useCallback(async () => {
    if (isDesktop) {
      if (!desktopSession) return
      const wf = await getWorkspaceFiles(desktopSession)
      setFiles(new Map(Object.entries(wf.files)))
      return
    }
    if (!dirHandle) return
    const loaded = await readWorkspaceFiles(dirHandle)
    setFiles(loaded)
  }, [dirHandle, desktopSession])

  // ブラウザモード: 承認済み operations を File System Access API で書き込む。
  const applyOperationsBrowser = useCallback(async (ops: Operation[]): Promise<string[]> => {
    if (!dirHandle) throw new Error('ワークスペースが attach されていません')
    const applied: string[] = []

    for (const op of ops) {
      if (
        op.type === 'create_file' ||
        op.type === 'update_file' ||
        op.type === 'upsert_fact'
      ) {
        await writeFile(dirHandle, op.target, op.change.after, op.entity_id)
        applied.push(op.target)
      } else if (op.type === 'mark_fact_status') {
        // 既存 fact の status のみを更新する（Go applier と同挙動）。
        if (!op.entity_id) throw new Error('mark_fact_status には entity_id が必要です')
        if (!op.fact_status_after) throw new Error('mark_fact_status には fact_status_after が必要です')
        await markFactStatus(dirHandle, op.target, op.entity_id, op.fact_status_after)
        applied.push(op.target)
      }
      // delete_file は UI で明示確認を経てから別フローで行う（MVP ではスキップ）
    }

    await refreshFiles()
    return applied
  }, [dirHandle, refreshFiles])

  // 承認済み patch を workspace へ適用する。
  // desktop: Go 経路（/apply-patch。session 束縛＋ResolveWithin。ADR-006 の正規経路）
  // browser: File System Access API（従来経路）
  const applyPatch = useCallback(async (patch: Patch): Promise<string[]> => {
    if (isDesktop) {
      if (!desktopSession) throw new Error('ワークスペースが attach されていません')
      // attachment との整合（/apply-patch は workspace_id 一致を要求）を確実にする。
      const bound: Patch = { ...patch, session_id: desktopSession, workspace_id: workspaceId }
      const result = await applyPatchViaServer(desktopSession, bound)
      if (result.failed_ops?.length) {
        throw new Error(`一部の操作が失敗しました: ${result.failed_ops.join(' / ')}`)
      }
      await refreshFiles()
      return [...new Set(patch.operations.map(op => op.target))]
    }
    return applyOperationsBrowser(patch.operations)
  }, [desktopSession, workspaceId, refreshFiles, applyOperationsBrowser])

  const attached = isDesktop ? desktopSession !== null : dirHandle !== null

  return (
    <WorkspaceContext.Provider value={{ dirHandle, workspaceId, files, attached, attach, refreshFiles, applyPatch }}>
      {children}
    </WorkspaceContext.Provider>
  )
}

// --- ファイル読み取りヘルパー ---

async function readWorkspaceFiles(
  dir: FileSystemDirectoryHandle,
  prefix = '',
  acc = new Map<string, string>(),
  depth = 0,
): Promise<Map<string, string>> {
  if (depth > 4) return acc
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  for await (const [name, handle] of (dir as any).entries() as AsyncIterable<[string, FileSystemHandle]>) {
    const path = prefix ? `${prefix}/${name}` : name
    if (handle.kind === 'file' && (name.endsWith('.yaml') || name.endsWith('.yml') || name.endsWith('.md'))) {
      const file = await (handle as FileSystemFileHandle).getFile()
      acc.set(path, await file.text())
    } else if (handle.kind === 'directory') {
      await readWorkspaceFiles(handle as FileSystemDirectoryHandle, path, acc, depth + 1)
    }
  }
  return acc
}

// --- ファイル書き込みヘルパー ---

async function writeFile(
  root: FileSystemDirectoryHandle,
  target: string,
  content: unknown,
  entityId?: string,
) {
  const parts = target.split('/')
  const fileName = parts.pop()!
  let dir = root
  for (const part of parts) {
    dir = await dir.getDirectoryHandle(part, { create: true })
  }

  let text: string
  if (typeof content === 'string') {
    // create_file / update_file: ファイル本文をそのまま書き込む。
    text = content
  } else if (content !== null && typeof content === 'object') {
    // upsert_fact: 正本形（fact 配列）として既存ファイルへ fact を upsert する。
    const fact = content as Record<string, unknown>
    const factId = (fact['fact_id'] as string) ?? entityId ?? `fact-${Date.now()}`
    let existingText = ''
    try {
      const existingFile = await dir.getFileHandle(fileName)
      existingText = await (await existingFile.getFile()).text()
    } catch {
      // ファイルが存在しない場合は空文字（新規配列から開始）
    }
    text = upsertFactYaml(existingText, fact, factId)
  } else {
    text = ''
  }

  const fileHandle = await dir.getFileHandle(fileName, { create: true })
  const writable = await fileHandle.createWritable()
  await writable.write(text)
  await writable.close()
}

// markFactStatus は既存の facts ファイル内の 1 fact の status を更新する。
// 対象ファイル / fact が無ければエラーにする（新規作成はしない）。
async function markFactStatus(
  root: FileSystemDirectoryHandle,
  target: string,
  factId: string,
  status: string,
) {
  const parts = target.split('/')
  const fileName = parts.pop()!
  let dir = root
  for (const part of parts) {
    dir = await dir.getDirectoryHandle(part)
  }
  const fileHandle = await dir.getFileHandle(fileName)
  const existingText = await (await fileHandle.getFile()).text()

  const text = markFactStatusYaml(existingText, factId, status)

  const writable = await fileHandle.createWritable()
  await writable.write(text)
  await writable.close()
}
