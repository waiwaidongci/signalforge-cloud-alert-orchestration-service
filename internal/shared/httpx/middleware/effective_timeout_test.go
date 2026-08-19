package middleware

import (
	"testing"
	"time"
)

func TestEffectiveTimeoutHonorsArgument(t *testing.T) {
	if got := effectiveTimeout(3 * time.Second); got != 3*time.Second {
		t.Fatalf("expected 3s, got %v", got)
	}
}
