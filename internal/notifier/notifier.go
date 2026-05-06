package notifier

import "fmt"

// Backend is the interface implemented by all notification backends.
type Backend interface {
	Send(subject, body string) error
}

// Notifier wraps a Backend and provides a unified Send method.
type Notifier struct {
	backend Backend
}

// New constructs a Notifier for the given backend type and options.
// Supported backends: "slack", "email", "webhook".
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
	default:
		return nil, fmt.Errorf("notifier: unknown backend %q", backendType)
	}
	if err != nil {
		return nil, err
	}
	return &Notifier{backend: b}, nil
}

// Send dispatches an alert via the configured backend.
// If no backend is configured it is a no-op.
func (n *Notifier) Send(subject, body string) error {
	if n == nil || n.backend == nil {
		return nil
	}
	return n.backend.Send(subject, body)
}
