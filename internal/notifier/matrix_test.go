package notifier

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMatrixBackend_MissingHomeserver(t *testing.T) {
	_, err := newMatrixBackend(map[string]string{
		"token":   "tok",
		"room_id": "!abc:example.com",
	})
	if err == nil || !strings.Contains(err.Error(), "homeserver") {
		t.Fatalf("expected homeserver error, got %v", err)
	}
}

func TestNewMatrixBackend_MissingToken(t *testing.T) {
	_, err := newMatrixBackend(map[string]string{
		"homeserver": "https://matrix.example.com",
		"room_id":    "!abc:example.com",
	})
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestNewMatrixBackend_MissingRoomID(t *testing.T) {
	_, err := newMatrixBackend(map[string]string{
		"homeserver": "https://matrix.example.com",
		"token":      "tok",
	})
	if err == nil || !strings.Contains(err.Error(), "room_id") {
		t.Fatalf("expected room_id error, got %v", err)
	}
}

func TestNewMatrixBackend_Valid(t *testing.T) {
	b, err := newMatrixBackend(map[string]string{
		"homeserver": "https://matrix.example.com",
		"token":      "tok",
		"room_id":    "!abc:example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestMatrixBackend_Send_Success(t *testing.T) {
	var gotAuth, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		buf := new(strings.Builder)
		buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b := &matrixBackend{
		homeserver: server.URL,
		token:      "secret",
		roomID:     "!room:example.com",
		client:     server.Client(),
	}
	if err := b.Send("Alert", "job failed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotAuth, "secret") {
		t.Errorf("expected bearer token in Authorization, got %q", gotAuth)
	}
	if !strings.Contains(gotBody, "Alert") || !strings.Contains(gotBody, "job failed") {
		t.Errorf("expected subject and body in payload, got %q", gotBody)
	}
}

func TestMatrixBackend_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	b := &matrixBackend{
		homeserver: server.URL,
		token:      "tok",
		roomID:     "!room:example.com",
		client:     server.Client(),
	}
	if err := b.Send("s", "b"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMatrixBackend_Send_BadURL(t *testing.T) {
	b := &matrixBackend{
		homeserver: "http://127.0.0.1:0",
		token:      "tok",
		roomID:     "!room:example.com",
		client:     &http.Client{},
	}
	if err := b.Send("s", "b"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
