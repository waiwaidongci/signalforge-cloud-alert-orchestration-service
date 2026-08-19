package db

import (
	"context"
	"testing"
)

func TestRunBatchHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := RunBatch(ctx, []int{1, 2}, func(context.Context, int) error { return nil }); err == nil {
		t.Fatal("expected cancelled context error")
	}
}
