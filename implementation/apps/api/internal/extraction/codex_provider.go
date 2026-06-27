package extraction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CodexExtractionProvider はユーザー自身の OpenAI アカウント（ADR-004）を使い、
// OpenAI Chat Completions API へ会話を投げて構造化 JSON を得る provider。
//
// 設計方針（Task4 Phase A）:
//   - AI 出力は untrusted JSON として扱い、解釈・正規化・schema 検証は
//     すべて呼び出し側のパイプライン（ExtractionService）に委ねる。
//   - credential（API key）は構築時に注入し、永続化・ログ出力・session/workspace
//     との混在を行わない。将来の request-scope 化を阻害しないよう、provider は
//     呼び出し側でリクエスト毎に構築できる軽量な値オブジェクトにする。
//   - Phase A はリトライを持たない単一試行。バックオフ等は Phase C。
type CodexExtractionProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

const (
	defaultCodexBaseURL = "https://api.openai.com/v1"
	defaultCodexModel   = "gpt-4o-mini"
	defaultCodexTimeout = 30 * time.Second
)

// CodexConfig は CodexExtractionProvider の構築パラメータ。
type CodexConfig struct {
	APIKey     string        // 必須。ユーザーの OpenAI API key（永続化しない）
	BaseURL    string        // 任意。既定は OpenAI 本番
	Model      string        // 任意。既定は defaultCodexModel
	Timeout    time.Duration // 任意。既定は defaultCodexTimeout
	HTTPClient *http.Client  // 任意。テスト用に注入可能
}

// NewCodexExtractionProvider は config から provider を構築する。
// API key が空の場合はエラー（サイレントに不正な provider を作らない）。
func NewCodexExtractionProvider(cfg CodexConfig) (*CodexExtractionProvider, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("codex provider requires an API key")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultCodexBaseURL
	}
	model := cfg.Model
	if model == "" {
		model = defaultCodexModel
	}
	client := cfg.HTTPClient
	if client == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = defaultCodexTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	return &CodexExtractionProvider{
		apiKey:  cfg.APIKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client:  client,
	}, nil
}

// Name は provider 名を返す。
func (c *CodexExtractionProvider) Name() string { return "codex" }

// --- OpenAI Chat Completions の最小リクエスト/レスポンス型 ---

type openAIChatRequest struct {
	Model          string                `json:"model"`
	Messages       []openAIChatMessage   `json:"messages"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatResponse struct {
	Choices []openAIChoice `json:"choices"`
}

type openAIChoice struct {
	Message openAIChatMessage `json:"message"`
}

// ExtractFacts は会話を OpenAI へ投げ、ExtractedFactResult を返す。
func (c *CodexExtractionProvider) ExtractFacts(ctx context.Context, conversation string) (*ExtractedFactResult, error) {
	payload := openAIChatRequest{
		Model: c.model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: getSystemPrompt()},
			{Role: "user", Content: getUserPrompt(conversation)},
		},
		// JSON mode を強制し、余計な散文やフェンスの混入を抑える。
		ResponseFormat: &openAIResponseFormat{Type: "json_object"},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("codex: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("codex: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("codex: request failed: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスは上限付きで読む（暴走防止）。
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("codex: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// 注意: API key はエラーに含めない（body と status のみ）。
		return nil, fmt.Errorf("codex: OpenAI returned status %d: %s", resp.StatusCode, truncate(string(body), 300))
	}

	var envelope openAIChatResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("codex: decode OpenAI envelope: %w", err)
	}
	if len(envelope.Choices) == 0 {
		return nil, fmt.Errorf("codex: OpenAI returned no choices")
	}

	content := stripCodeFence(envelope.Choices[0].Message.Content)

	var result ExtractedFactResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("codex: parse extraction JSON: %w", err)
	}
	return &result, nil
}

// stripCodeFence は ```json ... ``` や ``` ... ``` のコードフェンスを除去する。
// 実 LLM は JSON mode でもフェンスを混ぜることがあるため防御的に処理する。
func stripCodeFence(s string) string {
	t := strings.TrimSpace(s)
	if !strings.HasPrefix(t, "```") {
		return t
	}
	// 先頭フェンス行（```\n または ```json\n）を落とす。
	t = strings.TrimPrefix(t, "```")
	if i := strings.IndexByte(t, '\n'); i >= 0 {
		t = t[i+1:]
	} else {
		t = ""
	}
	// 末尾フェンスを落とす。
	t = strings.TrimSpace(t)
	t = strings.TrimSuffix(t, "```")
	return strings.TrimSpace(t)
}

// truncate はエラーメッセージ用に文字列をバイト長で切り詰める。
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
