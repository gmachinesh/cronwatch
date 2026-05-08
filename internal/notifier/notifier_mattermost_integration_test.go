package notifier_test

import (
	"testing"

	"github.com/cronwatch/cronwatch/internal/notifier"
)

func TestNew_MattermostBackend(t *testing.T) {
	n, err := notifier.New("mattermost", map[string]string{
		"webhook_url": "https://example.com/hooks/abc",
		"channel":     "#general",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNew_MattermostBackend_MissingOptions(t *testing.T) {
	_, err := notifier.New("mattermost", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}
