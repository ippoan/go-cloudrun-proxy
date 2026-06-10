package cloudrunproxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	h := HandleHealth("release-wave-gcp")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"ok":true,"service":"release-wave-gcp"}` {
		t.Fatalf("body = %q, want health JSON", got)
	}
}
