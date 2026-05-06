package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

type opsGenieBackend struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

func newOpsGenieBackend(opts map[string]string) (*opsGenieBackend, error) {
	apiKey, ok := opts["api_key"]
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("opsgenie: missing required option 'api_key'")
	}

	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = defaultOpsGenieURL
	}

	return &opsGenieBackend{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{},
	}, nil
}

func (b *opsGenieBackend) Send(subject, body string) error {
	payload := opsGeniePayload{
		Message:     subject,
		Description: body,
		Priority:    "P3",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, b.apiURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("opsgenie: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+b.apiKey)

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
