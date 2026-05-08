package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type mattermostBackend struct {
	webhookURL string
	channel    string
	username   string
	client     *http.Client
}

type mattermostPayload struct {
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	Text     string `json:"text"`
}

func newMattermostBackend(opts map[string]string) (*mattermostBackend, error) {
	url, ok := opts["webhook_url"]
	if !ok || url == "" {
		return nil, fmt.Errorf("mattermost: missing required option 'webhook_url'")
	}
	return &mattermostBackend{
		webhookURL: url,
		channel:    opts["channel"],
		username:   opts["username"],
		client:     &http.Client{},
	}, nil
}

func (b *mattermostBackend) Send(subject, body string) error {
	payload := mattermostPayload{
		Channel:  b.channel,
		Username: b.username,
		Text:     fmt.Sprintf("**%s**\n%s", subject, body),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}
	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("mattermost: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
