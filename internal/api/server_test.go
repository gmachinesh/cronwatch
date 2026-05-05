package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/internal/history"
	"github.com/cronwatch/internal/monitor"
	"github.com/cronwatch/internal/config"
)

func makeTestServer(t *testing.T) *Server {
	t.Helper()
	cfg := &config.Config{Jobs: []config.Job{
		{Name: "test-job", Schedule: "@every 1m", Command: "echo hi"},
	}}
	mon := monitor.New(cfg, nil)
	h, err := history.New(t.TempDir() + "/history.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	return New(":0", mon, h)
}

func TestHandleHealthz(t *testing.T) {
	s := makeTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	s.handleHealthz(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp healthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %q", resp.Status)
	}
}

func TestHandleHealthz_MethodNotAllowed(t *testing.T) {
	s := makeTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()
	s.handleHealthz(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleStatus_Empty(t *testing.T) {
	s := makeTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	s.handleStatus(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var entries []statusEntry
	if err := json.NewDecoder(w.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
}

func TestHandleHistory(t *testing.T) {
	s := makeTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/history?job=test-job", nil)
	w := httptest.NewRecorder()
	s.handleHistory(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
