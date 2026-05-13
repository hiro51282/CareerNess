import { useState, useRef, useEffect } from 'react'
import { sendMessage } from '../../api/client'
import { useWorkspace } from '../workspace/useWorkspace'
import type { Patch } from '../../types/patch'

interface Message {
  role: 'user' | 'ai'
  text: string
  patch?: Patch
}

interface Props {
  onPatchProposed: (patch: Patch) => void
}

export function ChatView({ onPatchProposed }: Props) {
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

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function handleSend() {
    const text = input.trim()
    if (!text || loading) return
    setInput('')
    setMessages(prev => [...prev, { role: 'user', text }])
    setLoading(true)

    try {
      const workspaceFiles: Record<string, string> = {}
      files.forEach((content, path) => { workspaceFiles[path] = content })

      const result = await sendMessage({ sessionId, workspaceId, message: text, workspaceFiles })
      setMessages(prev => [...prev, { role: 'ai', text: result.reply, patch: result.patch }])
      if (result.patch) {
        onPatchProposed(result.patch)
      }
    } catch (e) {
      setMessages(prev => [...prev, { role: 'ai', text: `エラーが発生しました: ${String(e)}` }])
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

  return (
    <div style={styles.container}>
      <div style={styles.messages}>
        {messages.map((m, i) => (
          <div key={i} style={{ ...styles.bubble, ...(m.role === 'user' ? styles.userBubble : styles.aiBubble) }}>
            <span style={styles.role}>{m.role === 'user' ? 'あなた' : 'AI'}</span>
            <p style={styles.text}>{m.text}</p>
            {m.patch && (
              <button
                style={styles.viewPatch}
                onClick={() => onPatchProposed(m.patch!)}
              >
                提案された変更を確認する →
              </button>
            )}
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
