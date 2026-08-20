package domain

import "context"

func ProbeCancelled(ctx context.Context, err error) bool {
	return err != nil && ctx.Err() != nil && err == ctx.Err()
}
