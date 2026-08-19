package middleware

import (
	"net/http"

	"github.com/acme/signalforge/internal/shared/httpx"
)

func PlaceholderAuth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := httpx.RequireToken(r, token); err != nil {
				httpx.WriteError(w, httpx.RequestID(r.Context()), err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
