package extraction

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// validExtractionJSON は OpenAI の message.content に入る妥当な抽出 JSON。
const validExtractionJSON = `{
  "extracted_facts": [
    {
      "type": "experience",
      "fact_id_hint": "payment-platform",
      "summary": "Payment platform migration",
      "period": {"start": "2022-01", "end": "2023-01"},
      "company": "ABC",
      "description": "Migrated payment platform to Go",
      "confidence": "high",
      "details": {},
      "extraction_notes": ["stated explicitly"],
      "clarification_questions": []
    }
  ],
  "extraction_quality": {
    "overall_confidence": "high",
    "completeness": "medium",
    "needs_clarification_count": 0,
    "summary": "ok"
  }
}`

// openAIEnvelope は content を OpenAI の chat completion レスポンス形に包む。
func openAIEnvelope(content string) string {
	b, _ := json.Marshal(openAIChatResponse{
		Choices: []openAIChoice{{Message: openAIChatMessage{Role: "assistant", Content: content}}},
	})
	return string(b)
}

// newProviderForServer は test server を指す provider を構築する。
func newProviderForServer(t *testing.T, url string) *CodexExtractionProvider {
	t.Helper()
	p, err := NewCodexExtractionProvider(CodexConfig{APIKey: "test-key", BaseURL: url})
	if err != nil {
		t.Fatalf("NewCodexExtractionProvider: %v", err)
	}
	return p
}

func TestCodex_Happy(t *testing.T) {
	var gotAuth, gotPath, gotMethod string
	var gotBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, openAIEnvelope(validExtractionJSON))
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	res, err := p.ExtractFacts(context.Background(), "2022年にABCで決済基盤をGoへ移行した")
	if err != nil {
		t.Fatalf("ExtractFacts error: %v", err)
	}

	// 解析結果
	if len(res.ExtractedFacts) != 1 {
		t.Fatalf("facts = %d, want 1", len(res.ExtractedFacts))
	}
	if res.ExtractedFacts[0].FactIDHint != "payment-platform" {
		t.Errorf("fact_id_hint = %q", res.ExtractedFacts[0].FactIDHint)
	}
	if res.ExtractionQuality.OverallConfidence != "high" {
		t.Errorf("overall_confidence = %q", res.ExtractionQuality.OverallConfidence)
	}

	// リクエストの妥当性
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/chat/completions" {
		t.Errorf("path = %q, want /chat/completions", gotPath)
	}
	if gotAuth != "Bearer test-key" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotBody["model"] == nil {
		t.Error("request body に model が無い")
	}
	if rf, ok := gotBody["response_format"].(map[string]any); !ok || rf["type"] != "json_object" {
		t.Errorf("response_format が json_object でない: %v", gotBody["response_format"])
	}
	msgs, ok := gotBody["messages"].([]any)
	if !ok || len(msgs) != 2 {
		t.Fatalf("messages = %v, want 2 件", gotBody["messages"])
	}
	// user メッセージに会話が埋め込まれていること
	last, _ := msgs[1].(map[string]any)
	if c, _ := last["content"].(string); !strings.Contains(c, "決済基盤") {
		t.Errorf("user メッセージに会話が含まれない: %q", c)
	}
}

func TestCodex_StripsFence(t *testing.T) {
	fenced := "```json\n" + validExtractionJSON + "\n```"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, openAIEnvelope(fenced))
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	res, err := p.ExtractFacts(context.Background(), "x")
	if err != nil {
		t.Fatalf("フェンス付き JSON の解析に失敗: %v", err)
	}
	if len(res.ExtractedFacts) != 1 {
		t.Fatalf("facts = %d, want 1", len(res.ExtractedFacts))
	}
}

func TestCodex_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, openAIEnvelope("これは JSON ではありません"))
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("不正な抽出 JSON はエラーになるべき")
	}
}

func TestCodex_Non200_SingleAttempt(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":{"message":"boom"}}`)
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("5xx はエラーになるべき")
	}
	// Phase A は単一試行（リトライしない）
	if n := calls.Load(); n != 1 {
		t.Errorf("呼び出し回数 = %d, want 1（Phase A はリトライしない）", n)
	}
}

func TestCodex_NoChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"choices":[]}`)
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	if _, err := p.ExtractFacts(context.Background(), "x"); err == nil {
		t.Fatal("choices 空はエラーになるべき")
	}
}

func TestCodex_ContextTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		io.WriteString(w, openAIEnvelope(validExtractionJSON))
	}))
	defer srv.Close()

	p := newProviderForServer(t, srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if _, err := p.ExtractFacts(ctx, "x"); err == nil {
		t.Fatal("ctx タイムアウトはエラーになるべき")
	}
}

func TestCodex_RequiresAPIKey(t *testing.T) {
	if _, err := NewCodexExtractionProvider(CodexConfig{APIKey: "  "}); err == nil {
		t.Fatal("API key 空は構築エラーになるべき")
	}
}

func TestStripCodeFence(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"plain", `{"a":1}`, `{"a":1}`},
		{"json fence", "```json\n{\"a\":1}\n```", `{"a":1}`},
		{"bare fence", "```\n{\"a\":1}\n```", `{"a":1}`},
		{"surrounding spaces", "  {\"a\":1}  ", `{"a":1}`},
	}
	for _, tc := range cases {
		if got := stripCodeFence(tc.in); got != tc.want {
			t.Errorf("%s: stripCodeFence(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}
