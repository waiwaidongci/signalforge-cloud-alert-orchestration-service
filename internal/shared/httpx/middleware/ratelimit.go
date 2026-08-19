package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	buckets map[string]*bucket
}

type bucket struct {
	windowStart time.Time
	count       int
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = 100
	}
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{limit: limit, window: window, buckets: make(map[string]*bucket)}
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := httpx.ClientIP(r)
		l.mu.Lock()
		now := time.Now()
		item := l.buckets[key]
		if item == nil || now.Sub(item.windowStart) >= l.window {
			item = &bucket{windowStart: now}
			l.buckets[key] = item
		}
		if item.count >= l.limit {
			l.mu.Unlock()
			httpx.WriteError(w, httpx.RequestID(r.Context()), apperr.New(apperr.KindRateLimited, "RATE_LIMITED", "请求过于频繁"))
			return
		}
		item.count++
		l.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
