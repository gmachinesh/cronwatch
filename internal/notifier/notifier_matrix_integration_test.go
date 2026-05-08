package notifier

import (
	"strings"
	"testing"
)

func TestNew_MatrixBackend(t *testing.T) {
	n, err := New("matrix", map[string]string{
		"homeserver": "https://matrix.example.com",
		"token":      "tok",
		"room_id":    "!abc:example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNew_MatrixBackend_MissingOptions(t *testing.T) {
	cases := []struct {
		name    string
		opts    map[string]string
		wantErr string
	}{
		{
			name:    "missing homeserver",
			opts:    map[string]string{"token": "tok", "room_id": "!r:x"},
			wantErr: "homeserver",
		},
		{
			name:    "missing token",
			opts:    map[string]string{"homeserver": "https://m.example.com", "room_id": "!r:x"},
			wantErr: "token",
		},
		{
			name:    "missing room_id",
			opts:    map[string]string{"homeserver": "https://m.example.com", "token": "tok"},
			wantErr: "room_id",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New("matrix", tc.opts)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}
