package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/acme/signalforge/internal/shared/apperr"
)

func MaxBodyBytes(maxBytes int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maxBytes > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func DecodeWithLimit(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return NormalizeBodyError(err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return NormalizeBodyError(err)
	}
	return nil
}

func ClientIP(r *http.Request) string {
	if value := r.Header.Get("X-Forwarded-For"); value != "" {
		return strings.Split(value, ",")[0]
	}
	if value := r.Header.Get("X-Real-IP"); value != "" {
		return value
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func IsAuthenticatedPlaceholder(token string) bool {
	return token != ""
}

func RequireToken(r *http.Request, token string) error {
	if token == "" {
		return nil
	}
	auth := r.Header.Get("Authorization")
	if auth != "Bearer "+token {
		return apperr.New(apperr.KindUnauthorized, "UNAUTHORIZED", "认证令牌无效")
	}
	return nil
}
