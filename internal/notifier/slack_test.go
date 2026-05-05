package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlackBackend_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b := newSlackBackend(server.URL)
	if err := b.Send("Test Alert", "Something failed"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSlackBackend_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	b := newSlackBackend(server.URL)
	err := b.Send("Test Alert", "Something failed")
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSlackBackend_Send_BadURL(t *testing.T) {
	b := newSlackBackend("http://127.0.0.1:0/nonexistent")
	err := b.Send("Test Alert", "body")
	if err == nil {
		t.Fatal("expected error for bad URL, got nil")
	}
}

func TestSlackBackend_Send_PayloadContainsSubjectAndBody(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		received = buf
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b := newSlackBackend(server.URL)
	if err := b.Send("MySubject", "MyBody"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := string(received)
	if len(payload) == 0 {
		t.Fatal("expected non-empty payload")
	}
	for _, want := range []string{"MySubject", "MyBody"} {
		if !containsStr(payload, want) {
			t.Errorf("payload missing %q: %s", want, payload)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && (
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}()))
}
