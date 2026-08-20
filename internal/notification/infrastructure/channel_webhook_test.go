package infrastructure

import (
	"context"
	"errors"
	"testing"
)

func TestWebhookHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := NewWebhookChannel().Send(ctx, "://bad-target", nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}
