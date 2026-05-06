package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOpsGenieBackend_MissingAPIKey(t *testing.T) {
	_, err := newOpsGenieBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing api_key, got nil")
	}
}

func TestNewOpsGenieBackend_Valid(t *testing.T) {
	b, err := newOpsGenieBackend(map[string]string{"api_key": "test-key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiKey != "test-key" {
		t.Errorf("expected api_key 'test-key', got %q", b.apiKey)
	}
	if b.apiURL != defaultOpsGenieURL {
		t.Errorf("expected default api_url, got %q", b.apiURL)
	}
}

func TestOpsGenieBackend_Send_Success(t *testing.T) {
	var received opsGeniePayload
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	b, _ := newOpsGenieBackend(map[string]string{"api_key": "my-key", "api_url": ts.URL})
	if err := b.Send("job failed", "details here"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Message != "job failed" {
		t.Errorf("expected message 'job failed', got %q", received.Message)
	}
	if received.Description != "details here" {
		t.Errorf("expected description 'details here', got %q", received.Description)
	}
	if authHeader != "GenieKey my-key" {
		t.Errorf("expected Authorization header 'GenieKey my-key', got %q", authHeader)
	}
}

func TestOpsGenieBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b, _ := newOpsGenieBackend(map[string]string{"api_key": "bad-key", "api_url": ts.URL})
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestOpsGenieBackend_Send_BadURL(t *testing.T) {
	b, _ := newOpsGenieBackend(map[string]string{"api_key": "k", "api_url": "http://127.0.0.1:0"})
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for bad URL, got nil")
	}
}
