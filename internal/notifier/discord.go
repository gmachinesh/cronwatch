package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type discordBackend struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content  string         `json:"content,omitempty"`
	Embeds   []discordEmbed `json:"embeds,omitempty"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

func newDiscordBackend(opts map[string]string) (*discordBackend, error) {
	url, ok := opts["webhook_url"]
	if !ok || url == "" {
		return nil, fmt.Errorf("discord notifier: missing required option 'webhook_url'")
	}
	return &discordBackend{
		webhookURL: url,
		client:     &http.Client{},
	}, nil
}

func (d *discordBackend) Send(subject, body string) error {
	payload := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:       subject,
				Description: body,
				Color:       15158332, // red
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord notifier: failed to marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("discord notifier: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
