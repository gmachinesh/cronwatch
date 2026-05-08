package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGotifyBackend_MissingURL(t *testing.T) {
	_, err := newGotifyBackend(map[string]string{"token": "abc123"})
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestNewGotifyBackend_MissingToken(t *testing.T) {
	_, err := newGotifyBackend(map[string]string{"url": "http://localhost"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewGotifyBackend_Valid(t *testing.T) {
	b, err := newGotifyBackend(map[string]string{
		"url":   "http://localhost:8080",
		"token": "mytoken",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.token != "mytoken" {
		t.Errorf("expected token 'mytoken', got %q", b.token)
	}
}

func TestGotifyBackend_Send_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, _ := newGotifyBackend(map[string]string{"url": ts.URL, "token": "tok"})
	if err := b.Send("Test Alert", "job failed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGotifyBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b, _ := newGotifyBackend(map[string]string{"url": ts.URL, "token": "bad"})
	if err := b.Send("Alert", "body"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGotifyBackend_Send_BadURL(t *testing.T) {
	b, _ := newGotifyBackend(map[string]string{"url": "http://127.0.0.1:0", "token": "tok"})
	if err := b.Send("Alert", "body"); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
