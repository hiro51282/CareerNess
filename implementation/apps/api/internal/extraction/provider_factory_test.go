package extraction

import "testing"

func TestNewProviderFromEnv_DefaultMock(t *testing.T) {
	// 空（未設定相当）は mock
	t.Setenv("EXTRACTION_PROVIDER", "")
	p, err := NewProviderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "mock" {
		t.Errorf("既定は mock であるべき, got %q", p.Name())
	}
}

func TestNewProviderFromEnv_ExplicitMock(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "mock")
	p, err := NewProviderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "mock" {
		t.Errorf("Name = %q, want mock", p.Name())
	}
}

func TestNewProviderFromEnv_Codex(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "codex")
	t.Setenv("OPENAI_API_KEY", "sk-test")
	p, err := NewProviderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "codex" {
		t.Errorf("Name = %q, want codex", p.Name())
	}
}

func TestNewProviderFromEnv_CodexMissingKey(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "codex")
	t.Setenv("OPENAI_API_KEY", "") // key 欠落
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("codex 指定で key 欠落は明示エラーになるべき（mock へフォールバックしない）")
	}
}

func TestNewProviderFromEnv_Unknown(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "bogus")
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("未知の provider 指定は明示エラーになるべき")
	}
}
