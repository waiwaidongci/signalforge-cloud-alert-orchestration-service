package app

import (
	"database/sql"
	"log/slog"

	alertpersistence "github.com/acme/signalforge/internal/alert/adapter/persistence"
	alertapp "github.com/acme/signalforge/internal/alert/application"
	dedupapp "github.com/acme/signalforge/internal/dedup/application"
	escalationpersistence "github.com/acme/signalforge/internal/escalation/adapter/persistence"
	escalationapp "github.com/acme/signalforge/internal/escalation/application"
	incidentpersistence "github.com/acme/signalforge/internal/incident/adapter/persistence"
	incidentapp "github.com/acme/signalforge/internal/incident/application"
	notificationpersistence "github.com/acme/signalforge/internal/notification/adapter/persistence"
	notificationapp "github.com/acme/signalforge/internal/notification/application"
	notificationinfra "github.com/acme/signalforge/internal/notification/infrastructure"
	routingpersistence "github.com/acme/signalforge/internal/routing/adapter/persistence"
	routingapp "github.com/acme/signalforge/internal/routing/application"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/metrics"
	silencepersistence "github.com/acme/signalforge/internal/silence/adapter/persistence"
	silenceapp "github.com/acme/signalforge/internal/silence/application"
	sourcepersistence "github.com/acme/signalforge/internal/source/adapter/persistence"
	sourceapp "github.com/acme/signalforge/internal/source/application"
)

type Dependencies struct {
	DB      *sql.DB
	Logger  *slog.Logger
	Clock   clock.Clock
	Metrics *metrics.Metrics

	Sources       *sourceapp.Service
	Alerts        *alertapp.Service
	Ingest        *alertapp.IngestService
	Incidents     *incidentapp.Service
	Silences      *silenceapp.Service
	Routing       *routingapp.Service
	Escalation    *escalationapp.Service
	Notifications *notificationapp.Service
}

func Build(raw *sql.DB, logger *slog.Logger, metrics *metrics.Metrics) (*Dependencies, error) {
	clk := clock.SystemClock{}

	sourceRepo := sourcepersistence.NewSQLRepository(raw)
	alertRepo := alertpersistence.NewSQLRepository(raw)
	incidentRepo := incidentpersistence.NewSQLRepository(raw)
	timelineRepo := incidentpersistence.NewTimelineSQLRepository(raw)
	silenceRepo := silencepersistence.NewSQLRepository(raw)
	routingRepo := routingpersistence.NewSQLRepository(raw)
	escalationRepo := escalationpersistence.NewSQLRepository(raw)
	notificationRepo := notificationpersistence.NewSQLRepository(raw)

	sources := sourceapp.NewService(sourceRepo, clk)
	alerts := alertapp.NewService(alertRepo, incidentRepo, timelineRepo, clk, logger, metrics)
	silences := silenceapp.NewService(silenceRepo, clk)
	routing := routingapp.NewService(routingRepo, clk)
	escalation := escalationapp.NewService(escalationRepo, clk)
	incidents := incidentapp.NewService(incidentRepo, timelineRepo, clk)

	dispatcher := notificationinfra.NewDispatcher(
		notificationinfra.NewLogChannel(logger),
		notificationinfra.NewWebhookChannel(),
		notificationinfra.NewEmailChannel(logger),
	)
	notifications := notificationapp.NewService(notificationRepo, dispatcher, clk)

	fingerprints := dedupapp.NewService(15 * 60 * 1_000_000_000)
	ingest := alertapp.NewIngestService(
		alertRepo, sourceRepo, fingerprints, silences, routing, notifications, incidentRepo, timelineRepo, clk, logger, metrics,
	)

	return &Dependencies{
		DB: raw, Logger: logger, Clock: clk, Metrics: metrics,
		Sources: sources, Alerts: alerts, Ingest: ingest, Incidents: incidents,
		Silences: silences, Routing: routing, Escalation: escalation, Notifications: notifications,
	}, nil
}
