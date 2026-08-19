package db

import "context"

func RunBatch[T any](ctx context.Context, items []T, fn func(context.Context, T) error) (result error) {
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(ctx, item); err != nil {
			return err
		}
	}
	return nil
}
