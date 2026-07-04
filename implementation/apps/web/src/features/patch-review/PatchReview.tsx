import { useState, useRef } from 'react'
import type { Patch, Operation } from '../../types/patch'
import { useWorkspace } from '../workspace/useWorkspace'

interface Props {
  patches: Patch[]
  // 適用結果の通知（ドメイン）。パネルを閉じる責務とは分離する。
  onApplied: (appliedPaths: string[]) => void
  // レビュー完了・中断時にパネルを閉じる（UI）。
  onClose: () => void
}

// 1 patch = 1 セマンティック変更（ADR-002）に沿い、複数 patch を逐次レビューする。
export function PatchReview({ patches, onApplied, onClose }: Props) {
  const { applyOperations } = useWorkspace()
  const [index, setIndex] = useState(0)
  const [applying, setApplying] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const appliedRef = useRef<string[]>([])

  const patch = patches[index]
  const isLast = index >= patches.length - 1
  const multi = patches.length > 1

  function finish() {
    onApplied(appliedRef.current)
    onClose()
  }

  // 現在の patch のレビューを終え、次へ進む（最後なら適用結果を通知して閉じる）。
  function advance() {
    if (isLast) finish()
    else setIndex(i => i + 1)
  }

  async function handleApprove() {
    setApplying(true)
    setError(null)
    try {
      const applied = await applyOperations(patch.operations)
      appliedRef.current = [...appliedRef.current, ...applied]
      advance()
    } catch (e) {
      setError(String(e))
    } finally {
      setApplying(false)
    }
  }

  const riskLevel = calcRisk(patch.operations)

  return (
    <div style={styles.overlay}>
      <div style={styles.panel}>
        <header style={styles.header}>
          <h2 style={styles.title}>
            変更提案を確認{multi ? `（${index + 1} / ${patches.length}）` : ''}
          </h2>
          <RiskBadge level={riskLevel} />
        </header>

        <p style={styles.summary}>{patch.summary}</p>
        <p style={styles.meta}>
          パッチ ID: {patch.patch_id} ／ {patch.operations.length} 件の変更
        </p>

        <div style={styles.opList}>
          {patch.operations.map(op => (
            <OperationCard key={op.op_id} op={op} />
          ))}
        </div>

        {error && <p style={styles.error}>{error}</p>}

        <div style={styles.actions}>
          <button style={styles.rejectBtn} onClick={advance} disabled={applying}>
            {isLast ? '却下して閉じる' : '却下して次へ'}
          </button>
          <button style={styles.approveBtn} onClick={handleApprove} disabled={applying}>
            {applying ? '適用中…' : isLast ? '承認して適用' : '承認して次へ'}
          </button>
        </div>
      </div>
    </div>
  )
}

function OperationCard({ op }: { op: Operation }) {
  const [open, setOpen] = useState(false)

  return (
    <div style={styles.opCard}>
      <div style={styles.opHeader} onClick={() => setOpen(v => !v)}>
        <span style={styles.opType}>{op.type}</span>
        <span style={styles.opTarget}>{op.target}</span>
        <ConfidenceBadge c={op.confidence} />
        <span style={styles.toggle}>{open ? '▲' : '▼'}</span>
      </div>
      {open && (
        <div style={styles.opDetail}>
          <p style={styles.rationale}>
            <strong>根拠：</strong>{op.rationale}
          </p>
          {op.entity_id && (
            <p style={styles.detailRow}>
              <strong>entity_id：</strong>{op.entity_id}
            </p>
          )}
          {op.fact_status_after && (
            <p style={styles.detailRow}>
              <strong>fact status：</strong>{op.fact_status_after}
            </p>
          )}
          <div style={styles.diffBlock}>
            <div style={styles.diffSide}>
              <p style={styles.diffLabel}>変更前</p>
              <pre style={styles.pre}>{JSON.stringify(op.change.before, null, 2) ?? 'null'}</pre>
            </div>
            <div style={styles.diffSide}>
              <p style={styles.diffLabel}>変更後</p>
              <pre style={styles.pre}>{JSON.stringify(op.change.after, null, 2)}</pre>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function RiskBadge({ level }: { level: 'high' | 'medium' | 'low' }) {
  const colors = { high: '#dc2626', medium: '#d97706', low: '#16a34a' }
  const labels = { high: '要注意', medium: '要確認', low: '低リスク' }
  return (
    <span style={{ ...styles.badge, background: colors[level] }}>
      {labels[level]}
    </span>
  )
}

function ConfidenceBadge({ c }: { c: Operation['confidence'] }) {
  const colors = { high: '#16a34a', medium: '#d97706', low: '#6b7280' }
  return (
    <span style={{ ...styles.confBadge, color: colors[c] }}>
      確信度: {c}
    </span>
  )
}

function calcRisk(ops: Operation[]): 'high' | 'medium' | 'low' {
  if (ops.some(op => op.type === 'delete_file')) return 'high'
  if (ops.some(op => op.type === 'upsert_fact' || op.type === 'mark_fact_status')) return 'medium'
  return 'low'
}

const styles: Record<string, React.CSSProperties> = {
  overlay: {
    position: 'fixed',
    inset: 0,
    background: 'rgba(0,0,0,0.4)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 100,
    padding: 16,
  },
  panel: {
    background: '#fff',
    borderRadius: 12,
    padding: 24,
    maxWidth: 700,
    width: '100%',
    maxHeight: '90vh',
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
    gap: 12,
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
    justifyContent: 'space-between',
  },
  title: { fontSize: 18, fontWeight: 700 },
  summary: { fontSize: 14, color: '#444', lineHeight: 1.6 },
  meta: { fontSize: 12, color: '#888' },
  opList: { display: 'flex', flexDirection: 'column', gap: 8 },
  opCard: { border: '1px solid #e5e5e5', borderRadius: 8, overflow: 'hidden' },
  opHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '8px 12px',
    cursor: 'pointer',
    background: '#f9f9f9',
    userSelect: 'none',
  },
  opType: {
    fontSize: 12,
    fontWeight: 600,
    background: '#1a1a1a',
    color: '#fff',
    padding: '2px 8px',
    borderRadius: 4,
  },
  opTarget: { fontSize: 13, flex: 1, color: '#333', fontFamily: 'monospace' },
  toggle: { fontSize: 11, color: '#999' },
  opDetail: { padding: 12, display: 'flex', flexDirection: 'column', gap: 8 },
  rationale: { fontSize: 13, color: '#444', lineHeight: 1.5 },
  detailRow: { fontSize: 12, color: '#555' },
  diffBlock: { display: 'flex', gap: 8 },
  diffSide: { flex: 1, minWidth: 0 },
  diffLabel: { fontSize: 11, fontWeight: 600, color: '#888', marginBottom: 4 },
  pre: {
    fontSize: 12,
    background: '#f5f5f5',
    padding: 8,
    borderRadius: 6,
    overflowX: 'auto',
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
  confBadge: { fontSize: 11, fontWeight: 600 },
  badge: {
    fontSize: 11,
    fontWeight: 700,
    color: '#fff',
    padding: '2px 8px',
    borderRadius: 20,
  },
  actions: { display: 'flex', gap: 8, justifyContent: 'flex-end', paddingTop: 8 },
  rejectBtn: {
    padding: '8px 20px',
    border: '1px solid #ddd',
    borderRadius: 8,
    background: '#fff',
    cursor: 'pointer',
    fontSize: 14,
  },
  approveBtn: {
    padding: '8px 20px',
    border: 'none',
    borderRadius: 8,
    background: '#1a1a1a',
    color: '#fff',
    cursor: 'pointer',
    fontSize: 14,
    fontWeight: 600,
  },
  error: { color: '#dc2626', fontSize: 13 },
}
