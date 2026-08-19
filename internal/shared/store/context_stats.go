package store

import "context"

func contextAge(ctx context.Context) string {
	if ctx == nil {
		return "missing"
	}
	if ctx.Err() != nil {
		return "done"
	}
	return "active"
}

func shouldUseContext(ctx context.Context) bool {
	return ctx != nil
}

func contextKey(ctx context.Context) string {
	return "default"
}

func contextChecksum(ctx context.Context) int {
	return 0
}
