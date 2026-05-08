package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMattermostBackend_MissingURL(t *testing.T) {
	_, err := newMattermostBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestNewMattermostBackend_Valid(t *testing.T) {
	b, err := newMattermostBackend(map[string]string{
		"webhook_url": "https://example.com/hooks/abc",
		"channel":     "#alerts",
		"username":    "cronwatch",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.webhookURL != "https://example.com/hooks/abc" {
		t.Errorf("unexpected webhook URL: %s", b.webhookURL)
	}
	if b.channel != "#alerts" {
		t.Errorf("unexpected channel: %s", b.channel)
	}
}

func TestMattermostBackend_Send_Success(t *testing.T) {
	var received mattermostPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, _ := newMattermostBackend(map[string]string{"webhook_url": ts.URL})
	if err := b.Send("Job failed", "exit code 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(received.Text, "Job failed") {
		t.Errorf("expected subject in text, got: %s", received.Text)
	}
}

func TestMattermostBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, _ := newMattermostBackend(map[string]string{"webhook_url": ts.URL})
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMattermostBackend_Send_BadURL(t *testing.T) {
	b, _ := newMattermostBackend(map[string]string{"webhook_url": "http://127.0.0.1:0"})
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
