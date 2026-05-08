package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type matrixBackend struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

type matrixMessage struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

func newMatrixBackend(opts map[string]string) (*matrixBackend, error) {
	homeserver := opts["homeserver"]
	if homeserver == "" {
		return nil, fmt.Errorf("matrix notifier: missing required option 'homeserver'")
	}
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("matrix notifier: missing required option 'token'")
	}
	roomID := opts["room_id"]
	if roomID == "" {
		return nil, fmt.Errorf("matrix notifier: missing required option 'room_id'")
	}
	return &matrixBackend{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}, nil
}

func (m *matrixBackend) Send(subject, body string) error {
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		m.homeserver, m.roomID)

	msg := matrixMessage{
		MsgType: "m.text",
		Body:    fmt.Sprintf("%s\n%s", subject, body),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("matrix notifier: failed to marshal message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("matrix notifier: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix notifier: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
