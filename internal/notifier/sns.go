package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type snsBackend struct {
	topicARN string
	region   string
	endpoint string // overridable for testing
	client   *http.Client
}

type snsPublishRequest struct {
	TopicARN string `json:"TopicArn"`
	Subject  string `json:"Subject"`
	Message  string `json:"Message"`
}

func newSNSBackend(opts map[string]string) (backend, error) {
	topicARN := opts["topic_arn"]
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic_arn is required")
	}
	region := opts["region"]
	if region == "" {
		return nil, fmt.Errorf("sns: region is required")
	}
	endpoint := opts["endpoint"]
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://sns.%s.amazonaws.com", region)
	}
	return &snsBackend{
		topicARN: topicARN,
		region:   region,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (s *snsBackend) Send(subject, body string) error {
	payload := snsPublishRequest{
		TopicARN: s.topicARN,
		Subject:  subject,
		Message:  body,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sns: marshal payload: %w", err)
	}
	resp, err := s.client.Post(s.endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("sns: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sns: unexpected status %d", resp.StatusCode)
	}
	return nil
}
