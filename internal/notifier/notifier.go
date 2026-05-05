package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertLevel represents the severity of an alert.
type AlertLevel string

const (
	AlertFailure AlertLevel = "failure"
	AlertDrift   AlertLevel = "drift"
)

// Alert holds the data for a single notification event.
type Alert struct {
	JobName   string     `json:"job_name"`
	Level     AlertLevel `json:"level"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
}

// Config holds notifier configuration.
type Config struct {
	SlackWebhookURL string
	HTTPTimeout     time.Duration
}

// Notifier dispatches alerts to configured backends.
type Notifier struct {
	cfg    Config
	client *http.Client
}

// New creates a new Notifier with the given configuration.
func New(cfg Config) *Notifier {
	timeout := cfg.HTTPTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Notifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

// Send dispatches an alert to all configured backends.
func (n *Notifier) Send(a Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}
	if n.cfg.SlackWebhookURL != "" {
		if err := n.sendSlack(a); err != nil {
			return fmt.Errorf("slack notify: %w", err)
		}
	}
	return nil
}

// sendSlack posts an alert message to a Slack incoming webhook.
func (n *Notifier) sendSlack(a Alert) error {
	text := fmt.Sprintf("[%s] *%s* — %s (at %s)",
		a.Level, a.JobName, a.Message, a.Timestamp.Format(time.RFC3339))

	payload := map[string]string{"text": text}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := n.client.Post(n.cfg.SlackWebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d from Slack", resp.StatusCode)
	}
	return nil
}
