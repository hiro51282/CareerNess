import { createContext, useState, useCallback, type ReactNode } from 'react'
import type { Operation } from '../../types/patch'
import * as yaml from 'js-yaml'

export interface WorkspaceContextType {
  dirHandle: FileSystemDirectoryHandle | null
  workspaceId: string
  files: Map<string, string>
  attach: () => Promise<void>
  refreshFiles: () => Promise<void>
  applyOperations: (ops: Operation[]) => Promise<string[]>
}

export const WorkspaceContext = createContext<WorkspaceContextType | null>(null)

export function WorkspaceProvider({ children }: { children: ReactNode }) {
  const [dirHandle, setDirHandle] = useState<FileSystemDirectoryHandle | null>(null)
  const [workspaceId, setWorkspaceId] = useState('local-careervault')
  const [files, setFiles] = useState<Map<string, string>>(new Map())

  const attach = useCallback(async () => {
    const handle = await window.showDirectoryPicker({ mode: 'readwrite' })
    console.debug(`[workspace] attached: ${handle.name}`)
    setDirHandle(handle)
    setWorkspaceId(handle.name)
    const loaded = await readWorkspaceFiles(handle)
    console.debug(`[workspace] files loaded: ${loaded.size}`, [...loaded.keys()])
    setFiles(loaded)
  }, [])

  const refreshFiles = useCallback(async () => {
    if (!dirHandle) return
    const loaded = await readWorkspaceFiles(dirHandle)
    setFiles(loaded)
  }, [dirHandle])

  // 承認済み operations をワークスペースに書き込む。apply は AI ではなくブラウザ側の責務。
  const applyOperations = useCallback(async (ops: Operation[]): Promise<string[]> => {
    if (!dirHandle) throw new Error('ワークスペースが attach されていません')
    const applied: string[] = []

    for (const op of ops) {
      if (
        op.type === 'create_file' ||
        op.type === 'update_file' ||
        op.type === 'upsert_fact'
      ) {
        await writeFile(dirHandle, op.target, op.change.after)
        applied.push(op.target)
      }
      // delete_file は UI で明示確認を経てから別フローで行う（MVP ではスキップ）
    }

    await refreshFiles()
    return applied
  }, [dirHandle, refreshFiles])

  return (
    <WorkspaceContext.Provider value={{ dirHandle, workspaceId, files, attach, refreshFiles, applyOperations }}>
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

async function writeFile(root: FileSystemDirectoryHandle, target: string, content: unknown) {
  const parts = target.split('/')
  const fileName = parts.pop()!
  let dir = root
  for (const part of parts) {
    dir = await dir.getDirectoryHandle(part, { create: true })
  }
  const fileHandle = await dir.getFileHandle(fileName, { create: true })
  const writable = await fileHandle.createWritable()

  let text: string
  if (typeof content === 'string') {
    text = content
  } else {
    let existing: Record<string, unknown> = {}
    try {
      const existingFile = await dir.getFileHandle(fileName)
      const existingContent = await (await existingFile.getFile()).text()
      existing = (yaml.load(existingContent) as Record<string, unknown>) ?? {}
    } catch {
      // ファイルが存在しない場合は空オブジェクト
    }

    if (typeof content === 'object' && content !== null && 'entity_id' in content) {
      // upsert: entity_id をキーにしてマージ
      const entity = content as Record<string, unknown>
      const id = entity['entity_id'] as string
      if (!existing['facts']) existing['facts'] = {}
      ;(existing['facts'] as Record<string, unknown>)[id] = entity
    } else {
      existing = { ...existing, ...(content as Record<string, unknown>) }
    }
    text = yaml.dump(existing, { lineWidth: 120 })
  }

  await writable.write(text)
  await writable.close()
}
