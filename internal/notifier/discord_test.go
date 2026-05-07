package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDiscordBackend_MissingURL(t *testing.T) {
	_, err := newDiscordBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestNewDiscordBackend_Valid(t *testing.T) {
	b, err := newDiscordBackend(map[string]string{"webhook_url": "https://discord.com/api/webhooks/test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestDiscordBackend_Send_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	b, _ := newDiscordBackend(map[string]string{"webhook_url": ts.URL})
	if err := b.Send("Test Alert", "Job failed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiscordBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, _ := newDiscordBackend(map[string]string{"webhook_url": ts.URL})
	if err := b.Send("Test Alert", "Job failed"); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDiscordBackend_Send_BadURL(t *testing.T) {
	b, _ := newDiscordBackend(map[string]string{"webhook_url": "http://127.0.0.1:0/invalid"})
	if err := b.Send("Test Alert", "Job failed"); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestDiscordBackend_Send_PayloadContainsSubjectAndBody(t *testing.T) {
	var captured []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		captured = buf
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	b, _ := newDiscordBackend(map[string]string{"webhook_url": ts.URL})
	b.Send("my-subject", "my-body")

	if !containsStr(string(captured), "my-subject") {
		t.Errorf("payload missing subject: %s", captured)
	}
	if !containsStr(string(captured), "my-body") {
		t.Errorf("payload missing body: %s", captured)
	}
}
