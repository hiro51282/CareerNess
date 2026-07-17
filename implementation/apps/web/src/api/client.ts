import type { MessageResponse, Patch, ChatTurn } from '../types/patch'

const BASE = '/api/v1'

export async function sendMessage(params: {
  sessionId: string
  workspaceId: string
  message: string
  history: ChatTurn[]
  workspaceFiles: Record<string, string>
}): Promise<MessageResponse> {
  const res = await fetch(`${BASE}/conversations/message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      session_id: params.sessionId,
      workspace_id: params.workspaceId,
      message: params.message,
      history: params.history,
      workspace_files: params.workspaceFiles,
    }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

// --- desktop モード（Go 経路）用 API。ブラウザモードでは使わない ---

export async function attachWorkspace(
  workspaceRoot: string,
): Promise<{ session_id: string; workspace_id: string }> {
  const res = await fetch(`${BASE}/workspace/attach`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ workspace_root: workspaceRoot }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export async function getWorkspaceFiles(
  sessionId: string,
): Promise<{ workspace_id: string; files: Record<string, string> }> {
  const res = await fetch(`${BASE}/workspace/files?session_id=${encodeURIComponent(sessionId)}`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export async function applyPatchViaServer(
  sessionId: string,
  patch: Patch,
): Promise<{ applied_count: number; failed_ops: string[] }> {
  const res = await fetch(`${BASE}/apply-patch`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ session_id: sessionId, patch }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

// AIStatus は AI 実行環境の利用可否（オンボーディング用）。
export interface AIStatus {
  provider: string
  ready: boolean
  detail: string
  guidance?: string
}

export async function getAIStatus(): Promise<AIStatus> {
  const res = await fetch(`${BASE}/ai/status`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function validatePatch(patch: Patch): Promise<void> {
  const res = await fetch(`${BASE}/patches/validate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(patch),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
}
