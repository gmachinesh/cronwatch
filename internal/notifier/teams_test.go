package notifier

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewTeamsBackend_MissingURL(t *testing.T) {
	_, err := newTeamsBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing webhook_url, got nil")
	}
	if !strings.Contains(err.Error(), "webhook_url") {
		t.Errorf("expected error to mention 'webhook_url', got: %v", err)
	}
}

func TestNewTeamsBackend_Valid(t *testing.T) {
	b, err := newTeamsBackend(map[string]string{"webhook_url": "https://example.com/hook"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.webhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook URL: %s", b.webhookURL)
	}
}

func TestTeamsBackend_Send_Success(t *testing.T) {
	var gotContentType string
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, _ := newTeamsBackend(map[string]string{"webhook_url": ts.URL})
	if err := b.Send("Job failed", "backup-db exited with code 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %s", gotContentType)
	}
	if !strings.Contains(gotBody, "Job failed") {
		t.Errorf("expected body to contain subject, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "backup-db exited with code 1") {
		t.Errorf("expected body to contain message body, got: %s", gotBody)
	}
}

func TestTeamsBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, _ := newTeamsBackend(map[string]string{"webhook_url": ts.URL})
	err := b.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got: %v", err)
	}
}

func TestTeamsBackend_Send_BadURL(t *testing.T) {
	b, _ := newTeamsBackend(map[string]string{"webhook_url": "http://127.0.0.1:0"})
	err := b.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for bad URL, got nil")
	}
}
