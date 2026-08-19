package application

import (
	"context"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
	"github.com/acme/signalforge/internal/silence/domain"
)

type fakeSilenceRepositoryForSuppression struct{}

func (fakeSilenceRepositoryForSuppression) CreateSilence(context.Context, domain.Silence) error { return nil }
func (fakeSilenceRepositoryForSuppression) UpdateSilence(context.Context, domain.Silence) error { return nil }
func (fakeSilenceRepositoryForSuppression) FindSilenceByID(context.Context, string) (domain.Silence, error) { return domain.Silence{}, nil }
func (fakeSilenceRepositoryForSuppression) ListSilences(context.Context, bool, int, int) ([]domain.Silence, int, error) { return nil, 0, nil }
func (fakeSilenceRepositoryForSuppression) DeleteSilence(context.Context, string) error { return nil }
func (fakeSilenceRepositoryForSuppression) ActiveSilences(context.Context, time.Time) ([]domain.Silence, error) { return nil, nil }
func (fakeSilenceRepositoryForSuppression) CreateSuppression(context.Context, domain.SuppressionRule) error { return nil }
func (fakeSilenceRepositoryForSuppression) UpdateSuppression(context.Context, domain.SuppressionRule) error { return nil }
func (fakeSilenceRepositoryForSuppression) FindSuppressionByID(context.Context, string) (domain.SuppressionRule, error) { return domain.SuppressionRule{}, nil }
func (fakeSilenceRepositoryForSuppression) ListSuppressions(context.Context, bool, int, int) ([]domain.SuppressionRule, int, error) { return nil, 0, nil }
func (fakeSilenceRepositoryForSuppression) DeleteSuppression(context.Context, string) error { return nil }
func (fakeSilenceRepositoryForSuppression) ActiveSuppressions(context.Context) ([]domain.SuppressionRule, error) { return nil, nil }

func TestIsSuppressedDoesNotMisdetectMissingSilence(t *testing.T) {
	service := NewService(fakeSilenceRepositoryForSuppression{}, clock.FixedClock{Time: time.Unix(0, 0).UTC()})
	suppressed, err := service.IsSuppressed(context.Background(), matcher.Target{}, severity.Info, false)
	if err != nil {
		t.Fatal(err)
	}
	if suppressed {
		t.Fatal("missing silence should not suppress")
	}
}
