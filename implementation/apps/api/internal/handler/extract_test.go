package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPostExtract_DefaultMock は、env 未設定（既定 mock）時に /extract が従来どおり
// patches を返すこと（A3 の配線で既定挙動が変わらないこと）を検証する。
func TestPostExtract_DefaultMock(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "") // 既定 = mock

	body := `{"conversation":"2022年にABCで決済基盤をGoへ移行した","session_id":"sess-x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostExtract(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "patches") {
		t.Errorf("レスポンスに patches が含まれない: %s", rec.Body.String())
	}
}

// TestPostExtract_NoFactsStill422 は、会話拡張（B）後も純抽出エンドポイント /extract は
// 抽出 0 件を 422 とする従来挙動を維持することを検証する（回帰）。
func TestPostExtract_NoFactsStill422(t *testing.T) {
	t.Setenv("EXTRACTION_PROVIDER", "") // 既定 = mock（疑問形 → 0 件）

	body := `{"conversation":"他にどんな情報を書けばいいですか？","session_id":"sess-x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extract", strings.NewReader(body))
	rec := httptest.NewRecorder()

	PostExtract(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422（/extract の 0 件挙動は不変）; body=%s", rec.Code, rec.Body.String())
	}
}
