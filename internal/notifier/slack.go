package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// slackPayload represents the JSON body sent to a Slack webhook.
type slackPayload struct {
	Text string `json:"text"`
}

// slackBackend sends alert messages to a Slack incoming webhook URL.
type slackBackend struct {
	webhookURL string
	client     *http.Client
}

func newSlackBackend(webhookURL string) *slackBackend {
	return &slackBackend{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *slackBackend) Send(subject, body string) error {
	text := fmt.Sprintf("*%s*\n%s", subject, body)
	payload := slackPayload{Text: text}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("slack: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}

	return nil
}
