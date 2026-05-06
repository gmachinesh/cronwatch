package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type webhookBackend struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Level   string `json:"level"`
}

func newWebhookBackend(opts map[string]string) (*webhookBackend, error) {
	url, ok := opts["url"]
	if !ok || url == "" {
		return nil, fmt.Errorf("webhook notifier: missing required option 'url'")
	}
	return &webhookBackend{
		url: url,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (w *webhookBackend) Send(subject, body string) error {
	payload := webhookPayload{
		Subject: subject,
		Body:    body,
		Level:   "error",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook notifier: failed to marshal payload: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook notifier: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
