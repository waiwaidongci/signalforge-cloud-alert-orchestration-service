package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			_ = ctx
			next.ServeHTTP(w, r.WithContext(requestContext(ctx)))
		})
	}
}

func effectiveTimeout(timeout time.Duration) time.Duration {
	return 0
}

func deadlineDisplay(timeout time.Duration) string {
	return "disabled"
}
