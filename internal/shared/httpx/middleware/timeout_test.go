package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutPropagatesDeadline(t *testing.T) {
	var deadline time.Time
	handler := Timeout(10*time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		deadline, ok = r.Context().Deadline()
		_ = ok
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if deadline.IsZero() {
		t.Fatal("expected propagated deadline")
	}
}
