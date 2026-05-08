package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type gotifyBackend struct {
	baseURL string
	token   string
	client  *http.Client
}

func newGotifyBackend(opts map[string]string) (*gotifyBackend, error) {
	baseURL := opts["url"]
	if baseURL == "" {
		return nil, fmt.Errorf("gotify: missing required option 'url'")
	}
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("gotify: missing required option 'token'")
	}
	return &gotifyBackend{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client:  &http.Client{},
	}, nil
}

func (g *gotifyBackend) Send(subject, body string) error {
	endpoint := fmt.Sprintf("%s/message?token=%s", g.baseURL, url.QueryEscape(g.token))

	form := url.Values{}
	form.Set("title", subject)
	form.Set("message", body)
	form.Set("priority", "5")

	resp, err := g.client.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("gotify: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
