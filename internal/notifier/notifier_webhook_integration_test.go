package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cronwatch/internal/notifier"
)

func TestNew_WebhookBackend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notifier.New("webhook", map[string]string{"url": ts.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Send("hello", "world"); err != nil {
		t.Fatalf("send failed: %v", err)
	}
}

func TestNew_WebhookBackend_MissingOptions(t *testing.T) {
	_, err := notifier.New("webhook", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing url option, got nil")
	}
}
