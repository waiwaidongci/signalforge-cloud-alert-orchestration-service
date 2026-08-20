package application

import (
	"context"

	"github.com/acme/signalforge/internal/source/infrastructure"
)

type ProbeRunner struct{}

func (ProbeRunner) Run(ctx context.Context) error {
	return infrastructure.ProbeContext(ctx)
}
