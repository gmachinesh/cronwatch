package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_SNSBackend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := New("sns", map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"region":    "us-east-1",
		"endpoint":  ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Send("alert subject", "alert body"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
}

func TestNew_SNSBackend_MissingOptions(t *testing.T) {
	_, err := New("sns", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing sns options")
	}
}
