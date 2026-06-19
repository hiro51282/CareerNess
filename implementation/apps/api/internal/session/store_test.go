package session

import "testing"

func TestStore_PutGet(t *testing.T) {
	s := NewStore()
	s.Put(Attachment{SessionID: "sess-1", WorkspaceID: "vault", WorkspaceRoot: "/home/u/vault"})

	got, ok := s.Get("sess-1")
	if !ok {
		t.Fatal("登録した session が取得できない")
	}
	if got.WorkspaceID != "vault" || got.WorkspaceRoot != "/home/u/vault" {
		t.Errorf("attachment 内容が不一致: %+v", got)
	}
}

func TestStore_GetMissing(t *testing.T) {
	s := NewStore()
	if _, ok := s.Get("unknown"); ok {
		t.Error("未登録 session は ok=false であるべき")
	}
}

func TestStore_PutOverwrites(t *testing.T) {
	s := NewStore()
	s.Put(Attachment{SessionID: "sess-1", WorkspaceID: "old", WorkspaceRoot: "/old"})
	s.Put(Attachment{SessionID: "sess-1", WorkspaceID: "new", WorkspaceRoot: "/new"})

	got, _ := s.Get("sess-1")
	if got.WorkspaceID != "new" || got.WorkspaceRoot != "/new" {
		t.Errorf("再 attach で上書きされていない: %+v", got)
	}
}
