package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookChannel struct {
	client *http.Client
}

func NewWebhookChannel() *WebhookChannel {
	return &WebhookChannel{client: &http.Client{Timeout: 5 * time.Second}}
}

func (c *WebhookChannel) Name() string {
	return "webhook"
}

func (c *WebhookChannel) Send(ctx context.Context, destination string, payload map[string]any) error {
	raw, _ := json.Marshal(payload)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, destination, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("build webhook request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", response.StatusCode)
	}
	return nil
}
