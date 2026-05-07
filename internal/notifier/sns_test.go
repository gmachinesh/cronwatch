package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSNSBackend_MissingTopicARN(t *testing.T) {
	_, err := newSNSBackend(map[string]string{"region": "us-east-1"})
	if err == nil {
		t.Fatal("expected error for missing topic_arn")
	}
}

func TestNewSNSBackend_MissingRegion(t *testing.T) {
	_, err := newSNSBackend(map[string]string{"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic"})
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestNewSNSBackend_Valid(t *testing.T) {
	b, err := newSNSBackend(map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"region":    "us-east-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestSNSBackend_Send_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b, err := newSNSBackend(map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"region":    "us-east-1",
		"endpoint":  ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Send("test subject", "test body"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
}

func TestSNSBackend_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b, err := newSNSBackend(map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"region":    "us-east-1",
		"endpoint":  ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSNSBackend_Send_BadURL(t *testing.T) {
	b, err := newSNSBackend(map[string]string{
		"topic_arn": "arn:aws:sns:us-east-1:123456789012:MyTopic",
		"region":    "us-east-1",
		"endpoint":  "http://127.0.0.1:0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Send("subject", "body"); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
