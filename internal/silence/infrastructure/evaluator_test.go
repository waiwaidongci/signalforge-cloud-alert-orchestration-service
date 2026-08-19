package infrastructure

import (
	"testing"
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
)

func TestMatchingSilenceReturnsNilWhenNoMatch(t *testing.T) {
	evaluator := NewEvaluator()
	if got := evaluator.MatchingSilence(nil, matcher.Target{}, time.Unix(0, 0).UTC()); got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}
