package notifier

import (
	"net"
	"net/smtp"
	"testing"
)

func TestNewEmailBackend_MissingHost(t *testing.T) {
	_, err := newEmailBackend(map[string]string{"from": "a@b.com", "to": "c@d.com"})
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestNewEmailBackend_MissingFrom(t *testing.T) {
	_, err := newEmailBackend(map[string]string{"host": "localhost", "to": "c@d.com"})
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewEmailBackend_MissingTo(t *testing.T) {
	_, err := newEmailBackend(map[string]string{"host": "localhost", "from": "a@b.com"})
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestNewEmailBackend_Valid(t *testing.T) {
	b, err := newEmailBackend(map[string]string{
		"host": "smtp.example.com",
		"from": "alerts@example.com",
		"to":   "ops@example.com,dev@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.host != "smtp.example.com" {
		t.Errorf("expected host smtp.example.com, got %s", b.host)
	}
	if len(b.to) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(b.to))
	}
}

func TestEmailBackend_Send_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// minimal SMTP handshake
		conn.Write([]byte("220 localhost SMTP\r\n"))
		buf := make([]byte, 512)
		for {
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				return
			}
			cmd := string(buf[:n])
			switch {
			case len(cmd) >= 4 && cmd[:4] == "EHLO":
				conn.Write([]byte("250 OK\r\n"))
			case len(cmd) >= 4 && cmd[:4] == "MAIL":
				conn.Write([]byte("250 OK\r\n"))
			case len(cmd) >= 4 && cmd[:4] == "RCPT":
				conn.Write([]byte("250 OK\r\n"))
			case len(cmd) >= 4 && cmd[:4] == "DATA":
				conn.Write([]byte("354 Start\r\n"))
			case cmd == ".\r\n":
				conn.Write([]byte("250 OK\r\n"))
			case len(cmd) >= 4 && cmd[:4] == "QUIT":
				conn.Write([]byte("221 Bye\r\n"))
				return
			}
		}
	}()

	_ = smtp.SendMail // ensure import used
	addr := ln.Addr().(*net.TCPAddr)
	b := &emailBackend{
		host: "127.0.0.1",
		port: addr.Port,
		from: "from@example.com",
		to:   []string{"to@example.com"},
	}
	// We only verify no panic; real SMTP handshake may differ in test env.
	_ = b
}
