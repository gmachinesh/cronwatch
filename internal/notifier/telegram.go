package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBase = "https://api.telegram.org/bot"

type telegramBackend struct {
	token  string
	chatID string
	client *http.Client
}

type telegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func newTelegramBackend(opts map[string]string) (*telegramBackend, error) {
	token, ok := opts["token"]
	if !ok || token == "" {
		return nil, fmt.Errorf("telegram notifier: missing required option 'token'")
	}
	chatID, ok := opts["chat_id"]
	if !ok || chatID == "" {
		return nil, fmt.Errorf("telegram notifier: missing required option 'chat_id'")
	}
	return &telegramBackend{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}, nil
}

func (t *telegramBackend) Send(subject, body string) error {
	text := fmt.Sprintf("*%s*\n%s", subject, body)
	msg := telegramMessage{
		ChatID: t.chatID,
		Text:   text,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("telegram notifier: failed to marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBase, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("telegram notifier: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
