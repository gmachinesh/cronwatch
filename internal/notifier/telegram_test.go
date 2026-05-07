package notifier

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewTelegramBackend_MissingToken(t *testing.T) {
	_, err := newTelegramBackend(map[string]string{"chat_id": "123"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Errorf("expected error to mention 'token', got: %v", err)
	}
}

func TestNewTelegramBackend_MissingChatID(t *testing.T) {
	_, err := newTelegramBackend(map[string]string{"token": "abc123"})
	if err == nil {
		t.Fatal("expected error for missing chat_id")
	}
	if !strings.Contains(err.Error(), "chat_id") {
		t.Errorf("expected error to mention 'chat_id', got: %v", err)
	}
}

func TestNewTelegramBackend_Valid(t *testing.T) {
	b, err := newTelegramBackend(map[string]string{"token": "abc123", "chat_id": "456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.token != "abc123" || b.chatID != "456" {
		t.Errorf("unexpected backend values: %+v", b)
	}
}

func TestTelegramBackend_Send_Success(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b := &telegramBackend{
		token:  "testtoken",
		chatID: "789",
		client: server.Client(),
	}
	// Override the API base via a custom URL by patching the Send method indirectly.
	// We test by pointing the client at a local server using a custom transport.
	b.client = &http.Client{
		Transport: &prefixRewriteTransport{prefix: telegramAPIBase + "testtoken", target: server.URL},
	}

	err := b.Send("Test Subject", "Test Body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "789") {
		t.Errorf("expected chat_id in payload, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "Test Subject") {
		t.Errorf("expected subject in payload, got: %s", gotBody)
	}
}

func TestTelegramBackend_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	b := &telegramBackend{
		token:  "testtoken",
		chatID: "789",
		client: &http.Client{
			Transport: &prefixRewriteTransport{prefix: telegramAPIBase + "testtoken", target: server.URL},
		},
	}
	err := b.Send("fail", "body")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected status 400 in error, got: %v", err)
	}
}

// prefixRewriteTransport rewrites requests whose URL starts with prefix to target.
type prefixRewriteTransport struct {
	prefix string
	target string
}

func (p *prefixRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	origURL := req.URL.String()
	newURL := strings.Replace(origURL, p.prefix, p.target, 1)
	newReq := req.Clone(req.Context())
	parsed, err := req.URL.Parse(newURL)
	if err != nil {
		return nil, err
	}
	newReq.URL = parsed
	newReq.Host = parsed.Host
	return http.DefaultTransport.RoundTrip(newReq)
}
