package db

import "context"

func RunBatch[T any](ctx context.Context, items []T, fn func(context.Context, T) error) (result error) {
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		defer func() { result = nil }()
		if result = fn(ctx, item); result != nil {
			return result
		}
	}
	return nil
}
