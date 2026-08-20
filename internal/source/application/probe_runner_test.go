package application

import (
	"context"
	"errors"
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
	"github.com/acme/signalforge/internal/source/infrastructure"
)

func TestProbeRunnerPropagatesRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := (ProbeRunner{}).Run(ctx)
	if !domain.ProbeCancelled(ctx, err) || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled probe must preserve request context error: %v", err)
	}
	if domain.ProbeCancelled(ctx, errors.New("other")) { t.Fatal("unrelated error must not be classified as cancellation") }
	if !errors.Is(infrastructure.ProbeContext(ctx), context.Canceled) { t.Fatal("infrastructure must observe cancellation directly") }
	if !infrastructure.ProbeContextsIndependent(ctx, context.Background()) { t.Fatal("request contexts must remain independent") }
	if err := (ProbeRunner{}).Run(context.Background()); err != nil {
		t.Fatalf("new probe must not inherit old cancellation: %v", err)
	}
}
