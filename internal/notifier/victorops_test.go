package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewVictorOpsBackend_MissingRoutingKey(t *testing.T) {
	_, err := newVictorOpsBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing routing_key")
	}
	if !strings.Contains(err.Error(), "routing_key") {
		t.Errorf("expected error to mention routing_key, got: %v", err)
	}
}

func TestNewVictorOpsBackend_Valid(t *testing.T) {
	b, err := newVictorOpsBackend(map[string]string{
		"routing_key": "test-key-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.routingKey != "test-key-123" {
		t.Errorf("expected routing_key 'test-key-123', got %q", b.routingKey)
	}
	if b.apiURL != defaultVictorOpsURL {
		t.Errorf("expected default API URL, got %q", b.apiURL)
	}
}

func TestVictorOpsBackend_Send_Success(t *testing.T) {
	var received victorOpsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, _ := newVictorOpsBackend(map[string]string{
		"routing_key": "my-route",
		"api_url":     ts.URL,
	})

	if err := b.Send("job failed", "exit code 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.EntityDisplayName != "job failed" {
		t.Errorf("expected entity_display_name 'job failed', got %q", received.EntityDisplayName)
	}
	if received.StateMessage != "exit code 1" {
		t.Errorf("expected state_message 'exit code 1', got %q", received.StateMessage)
	}
	if received.MessageType != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %q", received.MessageType)
	}
}

func TestVictorOpsBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, _ := newVictorOpsBackend(map[string]string{
		"routing_key": "my-route",
		"api_url":     ts.URL,
	})

	err := b.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got: %v", err)
	}
}

func TestVictorOpsBackend_Send_BadURL(t *testing.T) {
	b, _ := newVictorOpsBackend(map[string]string{
		"routing_key": "my-route",
		"api_url":     "http://127.0.0.1:0",
	})

	err := b.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
