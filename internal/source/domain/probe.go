package domain

import (
	"context"
	"errors"
)

func ProbeCancelled(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	ctxErr := ctx.Err()
	if ctxErr == nil {
		return false
	}
	return errors.Is(err, ctxErr)
}
