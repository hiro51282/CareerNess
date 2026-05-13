import { useState } from 'react'
import { WorkspaceProvider } from './features/workspace/WorkspaceContext'
import { useWorkspace } from './features/workspace/useWorkspace'
import { WorkspaceAttach } from './features/workspace/WorkspaceAttach'
import { ChatView } from './features/chat/ChatView'
import { PatchReview } from './features/patch-review/PatchReview'
import type { Patch } from './types/patch'

function AppContent() {
  const { dirHandle, workspaceId, files } = useWorkspace()
  const [pendingPatch, setPendingPatch] = useState<Patch | null>(null)
  const [appliedNotice, setAppliedNotice] = useState<string | null>(null)

  if (!dirHandle) {
    return <WorkspaceAttach />
  }

  function handleApproved(appliedPaths: string[]) {
    setPendingPatch(null)
    setAppliedNotice(`${appliedPaths.length} 件のファイルを更新しました: ${appliedPaths.join(', ')}`)
    setTimeout(() => setAppliedNotice(null), 4000)
  }

  return (
    <div style={styles.layout}>
      <aside style={styles.sidebar}>
        <div style={styles.sideHeader}>
          <span style={styles.vaultName}>📁 {workspaceId}</span>
          <span style={styles.fileCount}>{files.size} ファイル</span>
        </div>
        <div style={styles.fileList}>
          {Array.from(files.keys()).sort().map(path => (
            <div key={path} style={styles.fileItem} title={path}>
              {path}
            </div>
          ))}
          {files.size === 0 && (
            <p style={styles.empty}>ファイルがありません。<br/>facts/ などのフォルダを作成してください。</p>
          )}
        </div>
      </aside>

      <main style={styles.main}>
        {appliedNotice && (
          <div style={styles.notice}>{appliedNotice}</div>
        )}
        <ChatView onPatchProposed={setPendingPatch} />
      </main>

      {pendingPatch && (
        <PatchReview
          patch={pendingPatch}
          onApproved={handleApproved}
          onRejected={() => setPendingPatch(null)}
        />
      )}
    </div>
  )
}

export default function App() {
  return (
    <WorkspaceProvider>
      <AppContent />
    </WorkspaceProvider>
  )
}

const styles: Record<string, React.CSSProperties> = {
  layout: {
    display: 'flex',
    height: '100vh',
    overflow: 'hidden',
  },
  sidebar: {
    width: 220,
    background: '#fff',
    borderRight: '1px solid #e5e5e5',
    display: 'flex',
    flexDirection: 'column',
    flexShrink: 0,
  },
  sideHeader: {
    padding: '14px 12px 10px',
    borderBottom: '1px solid #e5e5e5',
    display: 'flex',
    flexDirection: 'column',
    gap: 2,
  },
  vaultName: {
    fontSize: 13,
    fontWeight: 700,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  fileCount: { fontSize: 11, color: '#888' },
  fileList: {
    flex: 1,
    overflowY: 'auto',
    padding: '8px 0',
  },
  fileItem: {
    padding: '4px 12px',
    fontSize: 12,
    color: '#444',
    fontFamily: 'monospace',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    cursor: 'default',
  },
  empty: {
    padding: '12px',
    fontSize: 12,
    color: '#aaa',
    lineHeight: 1.6,
  },
  main: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
    background: '#f8f9fa',
    position: 'relative',
  },
  notice: {
    position: 'absolute',
    top: 12,
    left: '50%',
    transform: 'translateX(-50%)',
    background: '#16a34a',
    color: '#fff',
    padding: '8px 16px',
    borderRadius: 8,
    fontSize: 13,
    zIndex: 10,
    whiteSpace: 'nowrap',
  },
}
