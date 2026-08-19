package application

import (
	"context"
	"testing"
	"time"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/severity"
)

type fakeAlertRepositoryForStatus struct {
	status alertdomain.Status
	alert  alertdomain.Alert
}

func (f *fakeAlertRepositoryForStatus) Create(context.Context, alertdomain.Alert) error { return nil }
func (f *fakeAlertRepositoryForStatus) Update(context.Context, alertdomain.Alert) error { return nil }
func (f *fakeAlertRepositoryForStatus) FindByID(context.Context, string) (alertdomain.Alert, error) {
	return f.alert, nil
}
func (f *fakeAlertRepositoryForStatus) FindBySourceExternalID(context.Context, string, string) (alertdomain.Alert, error) {
	return alertdomain.Alert{}, nil
}
func (f *fakeAlertRepositoryForStatus) FindByFingerprint(context.Context, string, int) ([]alertdomain.Alert, error) {
	return nil, nil
}
func (f *fakeAlertRepositoryForStatus) List(context.Context, alertdomain.ListFilter, int, int) ([]alertdomain.Alert, int, error) {
	return nil, 0, nil
}
func (f *fakeAlertRepositoryForStatus) UpdateIncidentID(context.Context, []string, string) error { return nil }
func (f *fakeAlertRepositoryForStatus) BatchUpdateStatus(_ context.Context, ids []string, status alertdomain.Status, at time.Time) error {
	f.status = status
	return nil
}
func (f *fakeAlertRepositoryForStatus) CountByIncident(context.Context, string, alertdomain.Status) (int, error) { return 0, nil }
func (f *fakeAlertRepositoryForStatus) HighestActiveSeverity(context.Context) (severity.Severity, error) { return severity.Info, nil }

type fakeTimelineForAlert struct{}

func (fakeTimelineForAlert) Append(context.Context, incidentdomain.TimelineEvent) error { return nil }

func TestAcknowledgeSetsAcknowledgedStatus(t *testing.T) {
	repo := &fakeAlertRepositoryForStatus{alert: alertdomain.Alert{ID: "alr_1", IncidentID: "inc_1"}}
	service := NewService(repo, nil, fakeTimelineForAlert{}, clock.FixedClock{Time: time.Unix(0, 0).UTC()}, nil, nil)
	if err := service.Acknowledge(context.Background(), []string{"alr_1"}, "operator"); err != nil {
		t.Fatal(err)
	}
	if repo.status != alertdomain.StatusAcknowledged {
		t.Fatalf("expected acknowledged, got %q", repo.status)
	}
}
