package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type teamsBackend struct {
	webhookURL string
	client     *http.Client
}

type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	ActivityText  string `json:"activityText"`
}

func newTeamsBackend(opts map[string]string) (*teamsBackend, error) {
	url, ok := opts["webhook_url"]
	if !ok || url == "" {
		return nil, fmt.Errorf("teams notifier: missing required option 'webhook_url'")
	}
	return &teamsBackend{
		webhookURL: url,
		client:     &http.Client{},
	}, nil
}

func (t *teamsBackend) Send(subject, body string) error {
	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "FF0000",
		Summary:    subject,
		Sections: []teamsSection{
			{
				ActivityTitle: subject,
				ActivityText:  body,
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams notifier: failed to marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("teams notifier: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("teams notifier: unexpected status code %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}
