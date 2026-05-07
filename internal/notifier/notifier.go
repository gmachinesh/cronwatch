package notifier

import "fmt"

// Notifier is the interface implemented by all alert backends.
type Notifier interface {
	Send(subject, body string) error
}

// New constructs a Notifier for the given backend type using the provided
// options map. Returns an error if the backend is unknown or misconfigured.
func New(backend string, opts map[string]string) (Notifier, error) {
	switch backend {
	case "slack":
		return newSlackBackend(opts)
	case "email":
		return newEmailBackend(opts)
	case "webhook":
		return newWebhookBackend(opts)
	case "pagerduty":
		return newPagerDutyBackend(opts)
	case "opsgenie":
		return newOpsGenieBackend(opts)
	case "teams":
		return newTeamsBackend(opts)
	case "victorops":
		return newVictorOpsBackend(opts)
	default:
		return nil, fmt.Errorf("notifier: unknown backend %q", backend)
	}
}
