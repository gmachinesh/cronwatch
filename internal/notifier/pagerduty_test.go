package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPagerDutyBackend_MissingRoutingKey(t *testing.T) {
	_, err := newPagerDutyBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing routing_key")
	}
}

func TestNewPagerDutyBackend_Valid(t *testing.T) {
	b, err := newPagerDutyBackend(map[string]string{"routing_key": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.routingKey != "abc123" {
		t.Errorf("expected routing_key abc123, got %s", b.routingKey)
	}
}

func TestPagerDutyBackend_Send_Success(t *testing.T) {
	var received pagerDutyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	b := &pagerDutyBackend{routingKey: "key1", client: ts.Client()}
	// Override URL by temporarily replacing the constant is not possible,
	// so we test via a real-ish call with a custom client pointing to test server.
	// Swap client transport to redirect to test server.
	b.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Host = ts.Listener.Addr().String()
			req.URL.Scheme = "http"
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	if err := b.Send("job failed", "exit code 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "key1" {
		t.Errorf("expected routing_key key1, got %s", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action trigger, got %s", received.EventAction)
	}
}

func TestPagerDutyBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := &pagerDutyBackend{
		routingKey: "key1",
		client: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				req.URL.Host = ts.Listener.Addr().String()
				req.URL.Scheme = "http"
				return http.DefaultTransport.RoundTrip(req)
			}),
		},
	}

	err := b.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
