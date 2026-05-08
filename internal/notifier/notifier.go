package notifier

import (
	"fmt"
)

// Backend is implemented by each notification provider.
type Backend interface {
	Send(subject, body string) error
}

// Notifier wraps a Backend and exposes a unified Send interface.
type Notifier struct {
	backend Backend
}

// New constructs a Notifier for the given backend type using the provided
// options map. Returns an error if the backend is unknown or misconfigured.
func New(backendType string, opts map[string]string) (*Notifier, error) {
	var (
		b   Backend
		err error
	)
	switch backendType {
	case "slack":
		b, err = newSlackBackend(opts)
	case "email":
		b, err = newEmailBackend(opts)
	case "webhook":
		b, err = newWebhookBackend(opts)
	case "pagerduty":
		b, err = newPagerDutyBackend(opts)
	case "opsgenie":
		b, err = newOpsGenieBackend(opts)
	case "teams":
		b, err = newTeamsBackend(opts)
	case "victorops":
		b, err = newVictorOpsBackend(opts)
	case "sns":
		b, err = newSNSBackend(opts)
	case "discord":
		b, err = newDiscordBackend(opts)
	case "telegram":
		b, err = newTelegramBackend(opts)
	case "gotify":
		b, err = newGotifyBackend(opts)
	case "mattermost":
		b, err = newMattermostBackend(opts)
	default:
		return nil, fmt.Errorf("notifier: unknown backend %q", backendType)
	}
	if err != nil {
		return nil, err
	}
	return &Notifier{backend: b}, nil
}

// Send delivers a notification with the given subject and body.
func (n *Notifier) Send(subject, body string) error {
	return n.backend.Send(subject, body)
}
