package httpx

import (
	"context"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func WithRequestID(r *http.Request, requestID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), requestIDKey, requestID))
}

func RequestID(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey).(string)
	if value == "" {
		return "unknown"
	}
	return value
}
