package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultVictorOpsURL = "https://alert.victorops.com/integrations/generic/20131114/alert"

type victorOpsBackend struct {
	routingKey string
	apiURL     string
	client     *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

func newVictorOpsBackend(opts map[string]string) (*victorOpsBackend, error) {
	rk, ok := opts["routing_key"]
	if !ok || rk == "" {
		return nil, fmt.Errorf("victorops: missing required option 'routing_key'")
	}

	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = defaultVictorOpsURL
	}

	return &victorOpsBackend{
		routingKey: rk,
		apiURL:     apiURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (v *victorOpsBackend) Send(subject, body string) error {
	payload := victorOpsPayload{
		MessageType:       "CRITICAL",
		EntityID:          subject,
		EntityDisplayName: subject,
		StateMessage:      body,
		Timestamp:         time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", v.apiURL, v.routingKey)
	resp, err := v.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("victorops: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
