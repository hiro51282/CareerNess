import { useState, useRef, useEffect, useCallback } from 'react'
import { sendMessage, getAIStatus, type AIStatus } from '../../api/client'
import { useWorkspace } from '../workspace/useWorkspace'
import type { Patch } from '../../types/patch'

interface Message {
  role: 'user' | 'ai'
  text: string
  patches?: Patch[]
  clarifications?: string[]
}

interface Props {
  onPatchesProposed: (patches: Patch[]) => void
}

export function ChatView({ onPatchesProposed }: Props) {
  const { workspaceId, files } = useWorkspace()
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'ai',
      text: 'こんにちは。キャリアについて教えてください。どんな仕事をしてきましたか？',
    },
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [sessionId] = useState(() => `sess-${Date.now()}`)
  const bottomRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const [aiStatus, setAiStatus] = useState<AIStatus | null>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  // AI 状態はハイブリッド方式: 起動時 1 回＋「再確認」ボタン＋送信エラー時に再取得。
  // 毎メッセージの事前チェックはしない（外部プロセス起動のコストと冗長性のため）。
  const refreshAIStatus = useCallback(() => {
    getAIStatus().then(setAiStatus).catch(() => setAiStatus(null))
  }, [])

  useEffect(() => {
    refreshAIStatus()
  }, [refreshAIStatus])

  async function handleSend() {
    const text = input.trim()
    if (!text || loading) return
    setInput('')
    setMessages(prev => [...prev, { role: 'user', text }])
    setLoading(true)

    try {
      const workspaceFiles: Record<string, string> = {}
      files.forEach((content, path) => { workspaceFiles[path] = content })

      // 直近の会話履歴を同送する（この時点の messages は今回の発言を含まない＝履歴そのもの）。
      // 上限はバックエンド側でも制限されるが、送信量も直近 10 ターンに抑える。
      const history = messages.slice(-10).map(m => ({ role: m.role, text: m.text }))

      const result = await sendMessage({ sessionId, workspaceId, message: text, history, workspaceFiles })
      setMessages(prev => [...prev, {
        role: 'ai',
        text: result.reply,
        patches: result.patches,
        clarifications: result.clarifications,
      }])
      if (result.patches?.length) {
        onPatchesProposed(result.patches)
      }
    } catch (e) {
      setMessages(prev => [...prev, { role: 'ai', text: `エラーが発生しました: ${String(e)}` }])
      // 認証切れ等で状態が変わった可能性があるため、AI 状態を再取得してバナーを最新化する。
      refreshAIStatus()
    } finally {
      setLoading(false)
    }
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  // 質問チップのタップで入力欄に Q/A テンプレをプリフィルし、回答を促す。
  // 回答は履歴付きの次ターンで送られ、同じ fact の enrich（C2）につながる。
  function handleClarificationTap(question: string) {
    setInput(`Q. ${question}\nA. `)
    textareaRef.current?.focus()
  }

  return (
    <div style={styles.container}>
      {aiStatus && (
        <div style={styles.statusBar}>
          <span style={{ ...styles.statusBadge, ...(aiStatus.ready ? styles.statusOk : styles.statusWarn) }}>
            AI: {aiStatus.provider === 'mock' ? 'Mock（開発用）' : aiStatus.provider}
            {aiStatus.ready ? '' : ' ⚠'}
          </span>
          {aiStatus.guidance && (
            <>
              <span style={styles.statusGuidance}>{aiStatus.guidance}</span>
              <button style={styles.statusRefresh} onClick={refreshAIStatus}>再確認</button>
            </>
          )}
        </div>
      )}
      <div style={styles.messages}>
        {messages.map((m, i) => (
          <div key={i} style={{ ...styles.bubble, ...(m.role === 'user' ? styles.userBubble : styles.aiBubble) }}>
            <span style={styles.role}>{m.role === 'user' ? 'あなた' : 'AI'}</span>
            <p style={styles.text}>{m.text}</p>
            {m.clarifications?.length ? (
              <div style={styles.chips}>
                {m.clarifications.map((q, qi) => (
                  <button
                    key={qi}
                    style={styles.chip}
                    onClick={() => handleClarificationTap(q)}
                    title="クリックすると入力欄に回答テンプレートが入ります"
                  >
                    💬 {q}
                  </button>
                ))}
              </div>
            ) : null}
            {m.patches?.length ? (
              <button
                style={styles.viewPatch}
                onClick={() => onPatchesProposed(m.patches!)}
              >
                提案された変更を確認する →
              </button>
            ) : null}
          </div>
        ))}
        {loading && (
          <div style={{ ...styles.bubble, ...styles.aiBubble }}>
            <span style={styles.role}>AI</span>
            <p style={styles.text}>考え中…</p>
          </div>
        )}
        <div ref={bottomRef} />
      </div>

      <div style={styles.inputArea}>
        <textarea
          ref={textareaRef}
          style={styles.textarea}
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="キャリアについて教えてください（Enter で送信、Shift+Enter で改行）"
          rows={3}
          disabled={loading}
        />
        <button style={styles.sendBtn} onClick={handleSend} disabled={loading || !input.trim()}>
          送信
        </button>
      </div>
    </div>
  )
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    gap: 0,
  },
  statusBar: {
    display: 'flex',
    alignItems: 'center',
    flexWrap: 'wrap',
    gap: 8,
    padding: '6px 20px',
    borderBottom: '1px solid #eee',
    background: '#fafafa',
    fontSize: 12,
    flexShrink: 0,
  },
  statusBadge: {
    fontWeight: 700,
    padding: '2px 8px',
    borderRadius: 10,
  },
  statusOk: { color: '#166534', background: '#dcfce7' },
  statusWarn: { color: '#92400e', background: '#fef3c7' },
  statusGuidance: { color: '#92400e', flex: 1, minWidth: 200, lineHeight: 1.5 },
  statusRefresh: {
    fontSize: 12,
    padding: '3px 10px',
    border: '1px solid #d97706',
    color: '#92400e',
    background: '#fff',
    borderRadius: 6,
    cursor: 'pointer',
  },
  messages: {
    flex: 1,
    overflowY: 'auto',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'column',
    gap: 12,
  },
  bubble: {
    maxWidth: 600,
    padding: '10px 14px',
    borderRadius: 10,
    display: 'flex',
    flexDirection: 'column',
    gap: 4,
  },
  userBubble: {
    alignSelf: 'flex-end',
    background: '#1a1a1a',
    color: '#fff',
  },
  aiBubble: {
    alignSelf: 'flex-start',
    background: '#fff',
    border: '1px solid #e5e5e5',
  },
  role: {
    fontSize: 11,
    fontWeight: 600,
    opacity: 0.6,
    textTransform: 'uppercase' as const,
    letterSpacing: '0.5px',
  },
  text: {
    fontSize: 14,
    lineHeight: 1.6,
    whiteSpace: 'pre-wrap' as const,
  },
  viewPatch: {
    alignSelf: 'flex-start',
    marginTop: 6,
    fontSize: 13,
    color: '#2563eb',
    background: 'none',
    border: 'none',
    cursor: 'pointer',
    padding: 0,
    textDecoration: 'underline',
  },
  chips: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: 6,
    marginTop: 6,
  },
  chip: {
    fontSize: 12,
    color: '#1d4ed8',
    background: '#eff6ff',
    border: '1px solid #bfdbfe',
    borderRadius: 14,
    padding: '4px 10px',
    cursor: 'pointer',
    textAlign: 'left',
    lineHeight: 1.4,
  },
  inputArea: {
    display: 'flex',
    gap: 8,
    padding: '12px 20px',
    borderTop: '1px solid #e5e5e5',
    background: '#fff',
  },
  textarea: {
    flex: 1,
    resize: 'none',
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: 8,
    fontSize: 14,
    fontFamily: 'inherit',
    outline: 'none',
  },
  sendBtn: {
    padding: '0 20px',
    background: '#1a1a1a',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    cursor: 'pointer',
    fontWeight: 600,
    fontSize: 14,
    opacity: 1,
  },
}
