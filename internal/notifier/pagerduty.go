package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

type pagerDutyBackend struct {
	routingKey string
	client     *http.Client
}

type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyInner    `json:"payload"`
}

type pagerDutyInner struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

func newPagerDutyBackend(opts map[string]string) (*pagerDutyBackend, error) {
	key, ok := opts["routing_key"]
	if !ok || key == "" {
		return nil, fmt.Errorf("pagerduty: missing required option 'routing_key'")
	}
	return &pagerDutyBackend{
		routingKey: key,
		client:     &http.Client{},
	}, nil
}

func (p *pagerDutyBackend) Send(subject, body string) error {
	payload := pagerDutyPayload{
		RoutingKey:  p.routingKey,
		EventAction: "trigger",
		Payload: pagerDutyInner{
			Summary:  subject + ": " + body,
			Source:   "cronwatch",
			Severity: "error",
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := p.client.Post(pagerDutyEventURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
