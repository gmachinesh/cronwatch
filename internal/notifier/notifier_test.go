package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cronwatch/internal/notifier"
)

func TestSend_NoBackend(t *testing.T) {
	n := notifier.New(notifier.Config{})
	err := n.Send(notifier.Alert{
		JobName: "backup",
		Level:   notifier.AlertFailure,
		Message: "exited with code 1",
	})
	if err != nil {
		t.Fatalf("expected no error with no backends configured, got: %v", err)
	}
}

func TestSend_SlackSuccess(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.New(notifier.Config{
		SlackWebhookURL: ts.URL,
		HTTPTimeout:     2 * time.Second,
	})

	a := notifier.Alert{
		JobName:   "db-backup",
		Level:     notifier.AlertDrift,
		Message:   "ran 5m late",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, ok := received["text"]
	if !ok {
		t.Fatal("slack payload missing 'text' field")
	}
	if !strings.Contains(text, "db-backup") {
		t.Errorf("expected job name in text, got: %s", text)
	}
	if !strings.Contains(text, string(notifier.AlertDrift)) {
		t.Errorf("expected alert level in text, got: %s", text)
	}
}

func TestSend_SlackNonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notifier.New(notifier.Config{
		SlackWebhookURL: ts.URL,
		HTTPTimeout:     2 * time.Second,
	})

	err := n.Send(notifier.Alert{
		JobName: "cleanup",
		Level:   notifier.AlertFailure,
		Message: "timed out",
	})
	if err == nil {
		t.Fatal("expected error on non-2xx response, got nil")
	}
}
