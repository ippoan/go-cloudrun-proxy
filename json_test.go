package cloudrunproxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, http.StatusCreated, map[string]any{"ok": true, "operation": "op-1"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("body is not valid JSON: %v (%q)", err, rec.Body.String())
	}
	if got["ok"] != true || got["operation"] != "op-1" {
		t.Fatalf("body = %v, want ok=true operation=op-1", got)
	}
}

func TestWriteJSONError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSONError(rec, http.StatusBadGateway, "cloud run upstream error")

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", rec.Code)
	}
	// body は固定文言のみ (upstream 詳細を echo しない規約の wire 形式を固定)。
	if got := strings.TrimSpace(rec.Body.String()); got != `{"error":"cloud run upstream error"}` {
		t.Fatalf("body = %q, want fixed error JSON", got)
	}
}
