package monitor

import (
	"fmt"
	"log"
	"time"
)

// AlertKind classifies the type of alert.
type AlertKind string

const (
	KindFailure AlertKind = "failure"
	KindDrift   AlertKind = "drift"
)

// Alert represents a notification event for a cron job.
type Alert struct {
	Job       string
	Kind      AlertKind
	Message   string
	Timestamp time.Time
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] job=%q kind=%s msg=%s", a.Timestamp.Format(time.RFC3339), a.Job, a.Kind, a.Message)
}

// AlertSender is the interface for dispatching alerts.
type AlertSender interface {
	Send(a Alert)
}

// LogAlertSender sends alerts to the standard logger.
type LogAlertSender struct{}

// Send logs the alert using the standard log package.
func (l *LogAlertSender) Send(a Alert) {
	a.Timestamp = time.Now()
	log.Printf("[alert] %s", a)
}
