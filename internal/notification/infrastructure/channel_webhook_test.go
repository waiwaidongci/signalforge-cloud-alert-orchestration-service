package infrastructure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookChannelHonorsCancelledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	channel := NewWebhookChannel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := channel.Send(ctx, server.URL, map[string]any{"title": "test"}); err == nil {
		t.Fatal("expected cancelled context error")
	}
}
