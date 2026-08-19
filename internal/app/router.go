package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	alerthttp "github.com/acme/signalforge/internal/alert/adapter/http"
	escalationhttp "github.com/acme/signalforge/internal/escalation/adapter/http"
	incidenthttp "github.com/acme/signalforge/internal/incident/adapter/http"
	notificationhttp "github.com/acme/signalforge/internal/notification/adapter/http"
	routinghttp "github.com/acme/signalforge/internal/routing/adapter/http"
	"github.com/acme/signalforge/internal/shared/httpx"
	"github.com/acme/signalforge/internal/shared/httpx/middleware"
	silencehttp "github.com/acme/signalforge/internal/silence/adapter/http"
	sourcehttp "github.com/acme/signalforge/internal/source/adapter/http"
)

type RouterOptions struct {
	AuthToken      string
	MaxBodyBytes   int64
	RequestTimeout time.Duration
	RateLimit      int
}

func NewRouter(deps *Dependencies, options RouterOptions) http.Handler {
	sourceHandler := sourcehttp.NewHandler(deps.Sources, deps.Logger)
	alertRest := alerthttp.NewRestHandler(deps.Alerts, deps.Logger)
	webhook := alerthttp.NewWebhookHandler(deps.Ingest, deps.Sources, deps.Logger)
	incidentHandler := incidenthttp.NewHandler(deps.Incidents, deps.Logger)
	silenceHandler := silencehttp.NewHandler(deps.Silences, deps.Logger)
	routingHandler := routinghttp.NewHandler(deps.Routing, deps.Logger)
	escalationHandler := escalationhttp.NewHandler(deps.Escalation, deps.Logger)
	notificationHandler := notificationhttp.NewHandler(deps.Notifications, deps.Logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}, httpx.RequestID(r.Context()))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.DB.PingContext(r.Context()); err != nil {
			httpx.WriteError(w, httpx.RequestID(r.Context()), err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"}, httpx.RequestID(r.Context()))
	})
	if deps.Metrics != nil {
		mux.Handle("GET /metrics", deps.Metrics.Handler())
	}

	mux.HandleFunc("GET /api/v1/sources", sourceHandler.List)
	mux.HandleFunc("POST /api/v1/sources", sourceHandler.Create)
	mux.HandleFunc("GET /api/v1/sources/{id}", sourceHandler.Get)
	mux.HandleFunc("PATCH /api/v1/sources/{id}", sourceHandler.Update)
	mux.HandleFunc("PUT /api/v1/sources/{id}", sourceHandler.Update)
	mux.HandleFunc("DELETE /api/v1/sources/{id}", sourceHandler.Delete)

	mux.HandleFunc("GET /api/v1/alerts", alertRest.List)
	mux.HandleFunc("GET /api/v1/alerts/{id}", alertRest.Get)
	mux.HandleFunc("POST /api/v1/alerts/acknowledge", alertRest.BatchAcknowledge)
	mux.HandleFunc("POST /api/v1/alerts/close", alertRest.BatchClose)
	mux.HandleFunc("POST /api/v1/webhook", webhook.ServeHTTP)

	mux.HandleFunc("GET /api/v1/incidents", incidentHandler.List)
	mux.HandleFunc("GET /api/v1/incidents/{id}", incidentHandler.Get)
	mux.HandleFunc("POST /api/v1/incidents/acknowledge", incidentHandler.Acknowledge)
	mux.HandleFunc("POST /api/v1/incidents/close", incidentHandler.Close)
	mux.HandleFunc("GET /api/v1/incidents/{id}/timeline", incidentHandler.Timeline)

	mux.HandleFunc("GET /api/v1/silences", silenceHandler.ListSilences)
	mux.HandleFunc("POST /api/v1/silences", silenceHandler.CreateSilence)
	mux.HandleFunc("GET /api/v1/silences/{id}", silenceHandler.GetSilence)
	mux.HandleFunc("PATCH /api/v1/silences/{id}", silenceHandler.UpdateSilence)
	mux.HandleFunc("PUT /api/v1/silences/{id}", silenceHandler.UpdateSilence)
	mux.HandleFunc("DELETE /api/v1/silences/{id}", silenceHandler.DeleteSilence)

	mux.HandleFunc("GET /api/v1/suppressions", silenceHandler.ListSuppressions)
	mux.HandleFunc("POST /api/v1/suppressions", silenceHandler.CreateSuppression)
	mux.HandleFunc("GET /api/v1/suppressions/{id}", silenceHandler.GetSuppression)
	mux.HandleFunc("PATCH /api/v1/suppressions/{id}", silenceHandler.UpdateSuppression)
	mux.HandleFunc("PUT /api/v1/suppressions/{id}", silenceHandler.UpdateSuppression)
	mux.HandleFunc("DELETE /api/v1/suppressions/{id}", silenceHandler.DeleteSuppression)

	mux.HandleFunc("GET /api/v1/routing-rules", routingHandler.List)
	mux.HandleFunc("POST /api/v1/routing-rules", routingHandler.Create)
	mux.HandleFunc("GET /api/v1/routing-rules/{id}", routingHandler.Get)
	mux.HandleFunc("PATCH /api/v1/routing-rules/{id}", routingHandler.Update)
	mux.HandleFunc("PUT /api/v1/routing-rules/{id}", routingHandler.Update)
	mux.HandleFunc("DELETE /api/v1/routing-rules/{id}", routingHandler.Delete)

	mux.HandleFunc("GET /api/v1/escalation-policies", escalationHandler.List)
	mux.HandleFunc("POST /api/v1/escalation-policies", escalationHandler.Create)
	mux.HandleFunc("GET /api/v1/escalation-policies/{id}", escalationHandler.Get)
	mux.HandleFunc("PATCH /api/v1/escalation-policies/{id}", escalationHandler.Update)
	mux.HandleFunc("PUT /api/v1/escalation-policies/{id}", escalationHandler.Update)
	mux.HandleFunc("DELETE /api/v1/escalation-policies/{id}", escalationHandler.Delete)

	mux.HandleFunc("GET /api/v1/notifications", notificationHandler.List)
	mux.HandleFunc("GET /api/v1/notifications/{id}", notificationHandler.Get)

	var handler http.Handler = mux
	handler = middleware.Recover(deps.Logger)(handler)
	handler = middleware.AccessLog(deps.Logger)(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.CORS(handler)
	handler = httpx.MaxBodyBytes(options.MaxBodyBytes, handler)
	handler = middleware.Timeout(options.RequestTimeout)(handler)
	handler = middleware.PlaceholderAuth(options.AuthToken)(handler)
	if options.RateLimit > 0 {
		handler = middleware.NewRateLimiter(options.RateLimit, time.Minute).Middleware(handler)
	}
	if deps.Metrics != nil {
		handler = observeHTTP(deps, handler)
	}
	return handler
}

type responseObserver struct {
	http.ResponseWriter
	status int
}

func (o *responseObserver) WriteHeader(status int) {
	o.status = status
	o.ResponseWriter.WriteHeader(status)
}

func observeHTTP(deps *Dependencies, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		observer := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(observer, r)
		if observer.status == 0 {
			observer.status = http.StatusOK
		}
		deps.Metrics.ObserveHTTP(r.Method, r.URL.Path, httpStatusText(observer.status), time.Since(start).Seconds())
	})
}

func httpStatusText(status int) string {
	if text := http.StatusText(status); text != "" {
		return text
	}
	return "unknown"
}

func healthJSON(db *sql.DB) []byte {
	raw, _ := json.Marshal(map[string]string{"status": "ok"})
	return raw
}
