import { useState } from 'react'
import { useWorkspace } from '../workspace/useWorkspace'
import type { Operation } from '../../types/patch'
import { collectFacts, type FactEntry, type FactRecord } from './collectFacts'

const TYPE_ORDER = ['experience', 'achievement', 'skill']
const TYPE_LABEL: Record<string, string> = {
  experience: '経験',
  achievement: '成果',
  skill: 'スキル',
  other: 'その他',
}

// FactList は attach 済みワークスペースの facts/*.yaml を読み取り専用で一覧表示する。
export function FactList() {
  const { files } = useWorkspace()
  const entries = collectFacts(files)

  if (entries.length === 0) {
    return (
      <div style={styles.empty}>
        まだ fact がありません。<br />
        チャットからキャリアについて話すと、fact 候補が提案されます。
      </div>
    )
  }

  // type ごとにグループ化（既知の型を優先順で、その他は後ろに）。
  const groups: Record<string, FactEntry[]> = {}
  for (const e of entries) {
    const t = typeof e.fact.type === 'string' ? e.fact.type : 'other'
    ;(groups[t] ||= []).push(e)
  }
  const sectionTypes = [
    ...TYPE_ORDER.filter((t) => groups[t]),
    ...Object.keys(groups).filter((t) => !TYPE_ORDER.includes(t)),
  ]

  return (
    <div style={styles.container}>
      <div style={styles.headerRow}>
        <h2 style={styles.heading}>Facts</h2>
        <span style={styles.total}>{entries.length} 件</span>
      </div>
      {sectionTypes.map((t) => (
        <section key={t} style={styles.section}>
          <h3 style={styles.sectionTitle}>
            {TYPE_LABEL[t] ?? t}
            <span style={styles.sectionCount}>{groups[t].length}</span>
          </h3>
          <div style={styles.cards}>
            {groups[t].map((e, i) => (
              <FactCard key={`${e.file}:${e.fact.fact_id ?? i}`} entry={e} />
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}

function FactCard({ entry }: { entry: FactEntry }) {
  const { applyOperations } = useWorkspace()
  const [open, setOpen] = useState(false)
  const [confirming, setConfirming] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const f: FactRecord = entry.fact
  const hasDetail = Boolean(f.description) || Boolean(f.company) || Boolean(f.period)

  // ユーザーが明示的に確定できるのは proposed / inferred かつ fact_id を持つ fact。
  const canConfirm = Boolean(f.fact_id) && (f.status === 'proposed' || f.status === 'inferred')

  async function handleConfirm() {
    setConfirming(true)
    setError(null)
    try {
      const op: Operation = {
        op_id: 'op-confirm',
        type: 'mark_fact_status',
        target: entry.file,
        entity_id: String(f.fact_id),
        change: { before: f.status ?? null, after: 'confirmed' },
        rationale: 'ユーザーが確定',
        confidence: 'high',
        fact_status_after: 'confirmed',
      }
      // 適用後 refreshFiles により status バッジが confirmed へ更新される。
      await applyOperations([op])
    } catch (e) {
      setError(String(e))
    } finally {
      setConfirming(false)
    }
  }

  return (
    <div style={styles.card}>
      <div style={styles.cardHead} onClick={() => hasDetail && setOpen((v) => !v)}>
        {f.status && <StatusBadge status={String(f.status)} />}
        <span style={styles.summary}>{f.summary || f.fact_id || '(no summary)'}</span>
        {f.confidence && <span style={styles.confidence}>確信度: {String(f.confidence)}</span>}
        {hasDetail && <span style={styles.toggle}>{open ? '▲' : '▼'}</span>}
      </div>

      {Array.isArray(f.tags) && f.tags.length > 0 && (
        <div style={styles.tags}>
          {f.tags.map((tag) => (
            <span key={String(tag)} style={styles.tag}>
              {String(tag)}
            </span>
          ))}
        </div>
      )}

      {open && hasDetail && (
        <div style={styles.detail}>
          {(f.company || f.period) && (
            <p style={styles.metaLine}>
              {f.company && <span>会社: {String(f.company)}</span>}
              {f.company && f.period && <span> ／ </span>}
              {f.period && <span>期間: {String(f.period)}</span>}
            </p>
          )}
          {f.description && <p style={styles.description}>{String(f.description)}</p>}
          <p style={styles.source}>
            {entry.file}
            {f.fact_id ? ` ／ ${String(f.fact_id)}` : ''}
          </p>
        </div>
      )}

      {(canConfirm || error) && (
        <div style={styles.footer}>
          {canConfirm && (
            <button style={styles.confirmBtn} onClick={handleConfirm} disabled={confirming}>
              {confirming ? '確定中…' : '確定にする'}
            </button>
          )}
          {error && <span style={styles.errorText}>{error}</span>}
        </div>
      )}
    </div>
  )
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    confirmed: '#16a34a',
    proposed: '#d97706',
    inferred: '#6b7280',
    rejected: '#dc2626',
  }
  return (
    <span style={{ ...styles.status, background: colors[status] ?? '#6b7280' }}>{status}</span>
  )
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, overflowY: 'auto', padding: '16px 20px' },
  headerRow: { display: 'flex', alignItems: 'baseline', gap: 10, marginBottom: 12 },
  heading: { fontSize: 18, fontWeight: 700 },
  total: { fontSize: 12, color: '#888' },
  section: { marginBottom: 20 },
  sectionTitle: {
    fontSize: 13,
    fontWeight: 700,
    color: '#555',
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    marginBottom: 8,
  },
  sectionCount: {
    fontSize: 11,
    color: '#888',
    background: '#eee',
    borderRadius: 10,
    padding: '1px 7px',
  },
  cards: { display: 'flex', flexDirection: 'column', gap: 8 },
  card: {
    background: '#fff',
    border: '1px solid #e5e5e5',
    borderRadius: 8,
    padding: '10px 12px',
    display: 'flex',
    flexDirection: 'column',
    gap: 6,
  },
  cardHead: { display: 'flex', alignItems: 'center', gap: 8, cursor: 'default' },
  status: {
    fontSize: 11,
    fontWeight: 700,
    color: '#fff',
    padding: '2px 8px',
    borderRadius: 4,
    textTransform: 'uppercase',
  },
  summary: { fontSize: 14, flex: 1, color: '#222', minWidth: 0 },
  confidence: { fontSize: 11, color: '#888', whiteSpace: 'nowrap' },
  toggle: { fontSize: 11, color: '#999' },
  tags: { display: 'flex', flexWrap: 'wrap', gap: 4 },
  tag: {
    fontSize: 11,
    color: '#2563eb',
    background: '#eff6ff',
    borderRadius: 4,
    padding: '1px 6px',
  },
  detail: {
    borderTop: '1px solid #f0f0f0',
    paddingTop: 8,
    display: 'flex',
    flexDirection: 'column',
    gap: 6,
  },
  metaLine: { fontSize: 12, color: '#555' },
  description: { fontSize: 13, color: '#444', lineHeight: 1.6, whiteSpace: 'pre-wrap' },
  source: { fontSize: 11, color: '#aaa', fontFamily: 'monospace' },
  footer: { display: 'flex', alignItems: 'center', gap: 8 },
  confirmBtn: {
    fontSize: 12,
    fontWeight: 600,
    color: '#fff',
    background: '#16a34a',
    border: 'none',
    borderRadius: 6,
    padding: '4px 12px',
    cursor: 'pointer',
    alignSelf: 'flex-start',
  },
  errorText: { fontSize: 12, color: '#dc2626' },
  empty: { padding: 32, fontSize: 14, color: '#888', lineHeight: 1.8, textAlign: 'center' },
}
