package cloudrunproxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testHeader = "X-Test-API-Key"
	testKey    = "secret-key-for-test-0123456789ab"
)

func authedRequest(t *testing.T, h http.Handler, key string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	if key != "" {
		req.Header.Set(testHeader, key)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRequireAPIKey(t *testing.T) {
	var called bool
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := RequireAPIKey(testHeader, testKey, next)

	t.Run("correct key passes", func(t *testing.T) {
		called = false
		rec := authedRequest(t, h, testKey)
		if rec.Code != http.StatusOK || !called {
			t.Fatalf("status = %d, called = %v; want 200, true", rec.Code, called)
		}
	})

	t.Run("wrong key rejected", func(t *testing.T) {
		called = false
		rec := authedRequest(t, h, "wrong-key")
		if rec.Code != http.StatusUnauthorized || called {
			t.Fatalf("status = %d, called = %v; want 401, false", rec.Code, called)
		}
	})

	t.Run("missing header rejected", func(t *testing.T) {
		called = false
		rec := authedRequest(t, h, "")
		if rec.Code != http.StatusUnauthorized || called {
			t.Fatalf("status = %d, called = %v; want 401, false", rec.Code, called)
		}
	})

	// 401 body は固定文言 JSON のみで、expected key を一切 echo しない
	// (値漏れ防止 regression)。
	t.Run("401 body does not echo key", func(t *testing.T) {
		rec := authedRequest(t, h, "wrong-key")
		body := rec.Body.String()
		if strings.Contains(body, testKey) {
			t.Fatalf("401 body echoes the expected key: %q", body)
		}
		if strings.TrimSpace(body) != `{"error":"unauthorized"}` {
			t.Fatalf("401 body = %q, want fixed unauthorized JSON", body)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", ct)
		}
	})

	// expected が空 (= secret 未注入) なら、header も空の "一致" でも
	// fail-closed で 401 (subtle.ConstantTimeCompare は空 vs 空で 1 を返す
	// ため、空チェックが無いと素通りする)。
	t.Run("empty expected fail-closed", func(t *testing.T) {
		called = false
		open := RequireAPIKey(testHeader, "", next)
		rec := authedRequest(t, open, "")
		if rec.Code != http.StatusUnauthorized || called {
			t.Fatalf("status = %d, called = %v; want 401, false", rec.Code, called)
		}
		rec = authedRequest(t, open, "anything")
		if rec.Code != http.StatusUnauthorized || called {
			t.Fatalf("status = %d, called = %v; want 401, false", rec.Code, called)
		}
	})
}
