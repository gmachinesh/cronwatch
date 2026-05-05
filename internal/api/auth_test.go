package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAPIKeyAuth_NoKeyConfigured(t *testing.T) {
	mw := apiKeyAuth("")(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 when no key configured, got %d", rec.Code)
	}
}

func TestAPIKeyAuth_MissingHeader(t *testing.T) {
	mw := apiKeyAuth("secret")(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", rec.Code)
	}
}

func TestAPIKeyAuth_WrongKey(t *testing.T) {
	mw := apiKeyAuth("secret")(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(authHeader, "wrong")
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong key, got %d", rec.Code)
	}
}

func TestAPIKeyAuth_CorrectKey(t *testing.T) {
	mw := apiKeyAuth("secret")(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(authHeader, "secret")
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for correct key, got %d", rec.Code)
	}
}

func TestAPIKeyAuth_EmptyProvidedKey(t *testing.T) {
	mw := apiKeyAuth("secret")(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(authHeader, "")
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)
	// empty string header is treated as missing
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for empty provided key, got %d", rec.Code)
	}
}
