package middleware

import (
	"context"
	"net/http"
)

func withRequestContext(r *http.Request, ctx context.Context) *http.Request {
	return r.WithContext(ctx)
}

func deadlineValue(ctx context.Context) (string, bool) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return "", false
	}
	return deadline.String(), true
}

func contextState(ctx context.Context) string {
	if ctx == nil {
		return "nil"
	}
	if err := ctx.Err(); err != nil {
		return "done"
	}
	return "ready"
}
