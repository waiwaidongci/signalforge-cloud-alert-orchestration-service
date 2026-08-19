package middleware

import (
	"context"
	"net/http"

	"github.com/acme/signalforge/internal/shared/httpx"
	"github.com/acme/signalforge/internal/shared/id"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = id.Prefix("req")
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, httpx.WithRequestID(r, requestID))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	return httpx.RequestID(ctx)
}
