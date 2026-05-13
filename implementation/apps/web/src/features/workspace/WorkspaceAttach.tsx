import { useWorkspace } from './useWorkspace'

export function WorkspaceAttach() {
  const { attach } = useWorkspace()

  async function handleClick() {
    try {
      await attach()
    } catch (e) {
      if (e instanceof DOMException && e.name === 'AbortError') return
      alert('ワークスペースの読み込みに失敗しました')
    }
  }

  return (
    <div style={styles.container}>
      <h1 style={styles.title}>CareerNess</h1>
      <p style={styles.desc}>
        キャリア情報を構造化・管理するためのローカルファーストワークスペース
      </p>
      <button style={styles.btn} onClick={handleClick}>
        CareerVault を開く
      </button>
      <p style={styles.note}>
        ※ Chrome / Edge が必要です。選択したフォルダの内容のみが使用されます。
      </p>
    </div>
  )
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100vh',
    gap: 16,
    padding: 24,
  },
  title: {
    fontSize: 32,
    fontWeight: 700,
    letterSpacing: '-0.5px',
  },
  desc: {
    color: '#555',
    maxWidth: 360,
    textAlign: 'center',
    lineHeight: 1.6,
  },
  btn: {
    padding: '12px 28px',
    fontSize: 16,
    fontWeight: 600,
    background: '#1a1a1a',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    cursor: 'pointer',
    marginTop: 8,
  },
  note: {
    fontSize: 12,
    color: '#888',
    maxWidth: 360,
    textAlign: 'center',
  },
}
