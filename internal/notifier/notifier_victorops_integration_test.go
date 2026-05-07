package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/cronwatch/internal/notifier"
)

func TestNew_VictorOpsBackend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notifier.New("victorops", map[string]string{
		"routing_key": "integration-key",
		"api_url":     ts.URL,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}

	if err := n.Send("test alert", "integration test body"); err != nil {
		t.Fatalf("Send() returned unexpected error: %v", err)
	}
}

func TestNew_VictorOpsBackend_MissingOptions(t *testing.T) {
	_, err := notifier.New("victorops", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing routing_key")
	}
}
