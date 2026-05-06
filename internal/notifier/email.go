package notifier

import (
	"fmt"
	"net/smtp"
	"strings"
)

type emailBackend struct {
	host     string
	port     int
	from     string
	to       []string
	username string
	password string
}

func newEmailBackend(cfg map[string]string) (*emailBackend, error) {
	host, ok := cfg["host"]
	if !ok || host == "" {
		return nil, fmt.Errorf("email notifier: missing 'host'")
	}
	from, ok := cfg["from"]
	if !ok || from == "" {
		return nil, fmt.Errorf("email notifier: missing 'from'")
	}
	toRaw, ok := cfg["to"]
	if !ok || toRaw == "" {
		return nil, fmt.Errorf("email notifier: missing 'to'")
	}
	port := 587
	return &emailBackend{
		host:     host,
		port:     port,
		from:     from,
		to:       strings.Split(toRaw, ","),
		username: cfg["username"],
		password: cfg["password"],
	}, nil
}

func (e *emailBackend) Send(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.from,
		strings.Join(e.to, ", "),
		subject,
		body,
	)
	var auth smtp.Auth
	if e.username != "" {
		auth = smtp.PlainAuth("", e.username, e.password, e.host)
	}
	return smtp.SendMail(addr, auth, e.from, e.to, []byte(msg))
}
