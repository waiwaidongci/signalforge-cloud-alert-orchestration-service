package application

import (
	"context"

	"github.com/acme/signalforge/internal/source/infrastructure"
)

type ProbeRunner struct{}

func cancelled(ctx context.Context) error { return ctx.Err() }

func probeContext(ctx context.Context) context.Context { return ctx }

func (ProbeRunner) Run(ctx context.Context) error {
	if err := cancelled(ctx); err != nil {
		return err
	}
	return infrastructure.ProbeContext(probeContext(ctx))
}
