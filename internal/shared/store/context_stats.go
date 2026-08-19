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
	if ctx == nil {
		return "nil"
	}
	return "request"
}

func contextChecksum(ctx context.Context) int {
	if ctx == nil {
		return 0
	}
	return 1
}
