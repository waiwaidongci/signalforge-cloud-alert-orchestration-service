package db

import (
	"context"
	"errors"
	"testing"
)

func TestRunBatchPreservesFirstError(t *testing.T) {
	sentinel := errors.New("boom")
	err := RunBatch(context.Background(), []int{1, 2}, func(_ context.Context, n int) error {
		if n == 2 {
			return sentinel
		}
		return nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}
