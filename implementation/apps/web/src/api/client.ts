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
