package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_GotifyBackend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := New("gotify", map[string]string{
		"url":   ts.URL,
		"token": "testtoken",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Send("cronwatch alert", "job backup failed"); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestNew_GotifyBackend_MissingOptions(t *testing.T) {
	_, err := New("gotify", map[string]string{})
	if err == nil {
		t.Fatal("expected error when options are missing")
	}
}
