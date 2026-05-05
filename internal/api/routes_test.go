package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes_HealthzRegistered(t *testing.T) {
	s := makeTestServer(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRoutes_StatusRegistered(t *testing.T) {
	s := makeTestServer(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRoutes_HistoryRegistered(t *testing.T) {
	s := makeTestServer(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	s.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRoutes_MethodNotAllowed(t *testing.T) {
	s := makeTestServer(t)

	for _, path := range []string{"/healthz", "/status", "/history"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, nil)
		s.Handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("path %s: expected 405, got %d", path, rec.Code)
		}
	}
}
