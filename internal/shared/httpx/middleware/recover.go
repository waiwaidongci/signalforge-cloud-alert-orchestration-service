package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/acme/signalforge/internal/shared/httpx"
)

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						"request_id", httpx.RequestID(r.Context()),
						"panic", recovered,
						"stack", string(debug.Stack()),
					)
					httpx.WriteError(w, httpx.RequestID(r.Context()), nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
