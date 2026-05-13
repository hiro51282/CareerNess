import { useContext } from 'react'
import { WorkspaceContext } from './WorkspaceContext'

export function useWorkspace() {
  const ctx = useContext(WorkspaceContext)
  if (!ctx) throw new Error('WorkspaceProvider の外で useWorkspace が呼ばれました')
  return ctx
}
