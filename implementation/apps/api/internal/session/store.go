// Package session は workspace attachment（session → workspace_root の束縛）を扱う。
//
// ADR-006 の最小実装：保持するのは transient なポインタ（root への参照）のみで、
// キャリアデータや workspace の mirror は持たない（ADR-003 / backend-structure）。
// User Session / Conversation / Patch Review の分離、失効、capability scope、
// revision による stale 判定は MVP では未実装で Task3 以降に持ち越す。
package session

import "sync"

// Attachment は 1 session に紐づく workspace 接続。MVP では最小フィールドのみ。
type Attachment struct {
	SessionID   string
	WorkspaceID string
	// WorkspaceRoot は正規化済み（絶対・symlink 解決済み）の実パス。
	// apply の封じ込め基準であり、リクエストボディからは受け取らず attach 時に確定する。
	WorkspaceRoot string
}

// Store は session_id → Attachment の in-memory 写像。
// truth ではなく揮発的な作業状態であり、プロセス内にのみ存在する。
type Store struct {
	mu   sync.RWMutex
	data map[string]Attachment
}

// NewStore は空の in-memory Store を生成する。
func NewStore() *Store {
	return &Store{data: make(map[string]Attachment)}
}

// Put は attachment を登録する。同一 session_id は上書きする（再 attach）。
func (s *Store) Put(a Attachment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[a.SessionID] = a
}

// Get は session_id に対応する attachment を返す。未登録なら ok=false。
func (s *Store) Get(sessionID string) (Attachment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.data[sessionID]
	return a, ok
}
