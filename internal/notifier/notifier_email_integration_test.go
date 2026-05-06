package notifier

import (
	"testing"
)

// TestNew_EmailBackend verifies that New wires up an email backend
// when the config specifies type=email with required fields.
func TestNew_EmailBackend(t *testing.T) {
	cfg := Config{
		Backend: "email",
		Options: map[string]string{
			"host": "smtp.example.com",
			"from": "alerts@example.com",
			"to":   "ops@example.com",
		},
	}
	n, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNew_EmailBackend_MissingOptions(t *testing.T) {
	cfg := Config{
		Backend: "email",
		Options: map[string]string{
			"host": "smtp.example.com",
			// missing from and to
		},
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing email options")
	}
}

func TestNew_UnknownBackend(t *testing.T) {
	cfg := Config{
		Backend: "carrier-pigeon",
		Options: map[string]string{},
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}
