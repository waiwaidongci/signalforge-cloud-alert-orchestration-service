package application

import (
	"context"
	"time"

	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
	"github.com/acme/signalforge/internal/shared/validator"
	"github.com/acme/signalforge/internal/silence/domain"
)

type Service struct {
	repository domain.Repository
	clock      clock.Clock
}

func NewService(repository domain.Repository, clk clock.Clock) *Service {
	if clk == nil {
		clk = clock.SystemClock{}
	}
	return &Service{repository: repository, clock: clk}
}

func (s *Service) CreateSilence(ctx context.Context, silence domain.Silence) (domain.Silence, error) {
	now := s.clock.Now()
	if err := validator.Run(
		validator.Required("name", silence.Name),
	); err != nil {
		return domain.Silence{}, err
	}
	if err := validator.TimeRange(silence.StartsAt, silence.EndsAt); err != nil {
		return domain.Silence{}, err
	}
	silence.ID = id.Prefix("sil")
	silence.CreatedAt = now
	silence.UpdatedAt = now
	if err := s.repository.CreateSilence(ctx, silence); err != nil {
		return domain.Silence{}, err
	}
	return silence, nil
}

func (s *Service) UpdateSilence(ctx context.Context, silence domain.Silence) (domain.Silence, error) {
	existing, err := s.repository.FindSilenceByID(ctx, silence.ID)
	if err != nil {
		return domain.Silence{}, err
	}
	if err := validator.TimeRange(silence.StartsAt, silence.EndsAt); err != nil {
		return domain.Silence{}, err
	}
	silence.CreatedAt = existing.CreatedAt
	silence.UpdatedAt = s.clock.Now()
	if err := s.repository.UpdateSilence(ctx, silence); err != nil {
		return domain.Silence{}, err
	}
	return silence, nil
}

func (s *Service) GetSilence(ctx context.Context, silenceID string) (domain.Silence, error) {
	return s.repository.FindSilenceByID(ctx, silenceID)
}

func (s *Service) ListSilences(ctx context.Context, activeOnly bool, limit, offset int) ([]domain.Silence, int, error) {
	return s.repository.ListSilences(ctx, activeOnly, time.Now().UTC(), limit, offset)
}

func (s *Service) DeleteSilence(ctx context.Context, silenceID string) error {
	return s.repository.DeleteSilence(ctx, silenceID)
}

func (s *Service) CreateSuppression(ctx context.Context, rule domain.SuppressionRule) (domain.SuppressionRule, error) {
	now := s.clock.Now()
	if err := validator.Run(
		validator.Required("name", rule.Name),
		validator.OneOf("high_severity", rule.HighSeverity.String(), severity.All()...),
		validator.OneOf("low_severity", rule.LowSeverity.String(), severity.All()...),
	); err != nil {
		return domain.SuppressionRule{}, err
	}
	rule.ID = id.Prefix("sup")
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.repository.CreateSuppression(ctx, rule); err != nil {
		return domain.SuppressionRule{}, err
	}
	return rule, nil
}

func (s *Service) UpdateSuppression(ctx context.Context, rule domain.SuppressionRule) (domain.SuppressionRule, error) {
	existing, err := s.repository.FindSuppressionByID(ctx, rule.ID)
	if err != nil {
		return domain.SuppressionRule{}, err
	}
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = s.clock.Now()
	if err := s.repository.UpdateSuppression(ctx, rule); err != nil {
		return domain.SuppressionRule{}, err
	}
	return rule, nil
}

func (s *Service) GetSuppression(ctx context.Context, suppressionID string) (domain.SuppressionRule, error) {
	return s.repository.FindSuppressionByID(ctx, suppressionID)
}

func (s *Service) ListSuppressions(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.SuppressionRule, int, error) {
	return s.repository.ListSuppressions(ctx, enabledOnly, limit, offset)
}

func (s *Service) DeleteSuppression(ctx context.Context, suppressionID string) error {
	return s.repository.DeleteSuppression(ctx, suppressionID)
}

func (s *Service) ActiveSilences(ctx context.Context) ([]domain.Silence, error) {
	return s.repository.ActiveSilences(ctx, s.clock.Now())
}

func (s *Service) IsSuppressed(ctx context.Context, target matcher.Target, currentSeverity severity.Severity, highActive bool) (bool, error) {
	if highActive {
		rules, err := s.repository.ActiveSuppressions(ctx)
		if err != nil {
			return false, err
		}
		for _, rule := range rules {
			if !rule.Enabled {
				continue
			}
			if rule.SourceMatcher.Matches(target) &&
				severity.HigherOrEqual(rule.HighSeverity, currentSeverity) {
				return true, nil
			}
		}
	}
	return false, nil
}
