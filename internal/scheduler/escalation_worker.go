package scheduler

import (
	"context"
	"time"

	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	notificationdomain "github.com/acme/signalforge/internal/notification/domain"
	"github.com/acme/signalforge/internal/shared/matcher"
)

func (r *Runner) runEscalation(ctx context.Context) {
	incidents, _, err := r.deps.Incidents.List(ctx, incidentdomain.ListFilter{Status: incidentdomain.StatusActive}, 100, 0)
	if err != nil {
		r.logger.Warn("list active incidents failed", "error", err)
		return
	}
	for _, incident := range incidents {
		policies, err := r.deps.Escalation.Match(ctx, matcher.Target{
			SourceID: incident.SourceID,
			Service:  incident.Title,
			Severity: incident.Severity,
		})
		if err != nil {
			r.logger.Warn("match escalation policy failed", "incident_id", incident.ID, "error", err)
			continue
		}
		for _, policy := range policies {
			if !policy.AllowsRepeat(0) {
				continue
			}
			if time.Since(incident.LastSeenAt) < policy.NextDelay(0) {
				continue
			}
			for _, route := range policy.Routes {
				notification := notificationdomain.Notification{
					IncidentID: incident.ID, Channel: route.Channel, Destination: route.Destination,
					Payload: map[string]any{"title": incident.Title, "severity": incident.Severity.String(), "event": "escalation"},
				}
				if err := r.deps.Notifications.Notify(ctx, notification); err != nil {
					r.logger.Warn("escalation notification failed", "incident_id", incident.ID, "channel", route.Channel, "error", err)
				}
			}
		}
	}
}
