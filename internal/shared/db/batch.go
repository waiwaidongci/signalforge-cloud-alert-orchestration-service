package db

import "context"

func RunBatch[T any](ctx context.Context, items []T, fn func(context.Context, T) error) error {
	if len(items) == 0 {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	for i := range items {
		if err := fn(ctx, items[i]); err != nil {
			return err
		}
	}
	return nil
}
