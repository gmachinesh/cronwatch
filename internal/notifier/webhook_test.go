package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWebhookBackend_MissingURL(t *testing.T) {
	_, err := newWebhookBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing url, got nil")
	}
}

func TestNewWebhookBackend_Valid(t *testing.T) {
	b, err := newWebhookBackend(map[string]string{"url": "http://example.com/hook"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.url != "http://example.com/hook" {
		t.Errorf("expected url to be set, got %q", b.url)
	}
}

func TestWebhookBackend_Send_Success(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, _ := newWebhookBackend(map[string]string{"url": ts.URL})
	if err := b.Send("Test Subject", "Test Body"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got %q", received.Subject)
	}
	if received.Body != "Test Body" {
		t.Errorf("expected body 'Test Body', got %q", received.Body)
	}
}

func TestWebhookBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, _ := newWebhookBackend(map[string]string{"url": ts.URL})
	if err := b.Send("subj", "body"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestWebhookBackend_Send_BadURL(t *testing.T) {
	b, _ := newWebhookBackend(map[string]string{"url": "http://127.0.0.1:0/nope"})
	if err := b.Send("subj", "body"); err == nil {
		t.Fatal("expected error for bad URL, got nil")
	}
}
