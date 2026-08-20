package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry      *prometheus.Registry
	mu            sync.RWMutex
	recentPaths   map[string]int
	HTTPRequests  *prometheus.CounterVec
	HTTPDuration  *prometheus.HistogramVec
	Alerts        *prometheus.CounterVec
	Notifications *prometheus.CounterVec
	SchedulerRuns prometheus.Counter
}

func New() *Metrics {
	registry := prometheus.NewRegistry()
	m := &Metrics{
		registry:    registry,
		recentPaths: make(map[string]int),
		HTTPRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "signalforge_http_requests_total",
			Help: "Total HTTP requests by method, path and status.",
		}, []string{"method", "path", "status"}),
		HTTPDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "signalforge_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),
		Alerts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "signalforge_alerts_total",
			Help: "Alerts received by severity and status.",
		}, []string{"severity", "status"}),
		Notifications: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "signalforge_notifications_total",
			Help: "Notifications sent by channel and status.",
		}, []string{"channel", "status"}),
		SchedulerRuns: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "signalforge_scheduler_runs_total",
			Help: "Total scheduler loop runs.",
		}),
	}
	registry.MustRegister(m.HTTPRequests, m.HTTPDuration, m.Alerts, m.Notifications, m.SchedulerRuns)
	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) ObserveHTTP(method, path, status string, seconds float64) {
	if m == nil {
		return
	}
	m.HTTPRequests.WithLabelValues(method, path, status).Inc()
	m.HTTPDuration.WithLabelValues(method, path).Observe(seconds)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recentPaths[path]++
}
