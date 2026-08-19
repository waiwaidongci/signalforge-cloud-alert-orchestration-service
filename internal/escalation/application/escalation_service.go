package application

import (
	"context"

	"github.com/acme/signalforge/internal/escalation/domain"
	"github.com/acme/signalforge/internal/shared/clock"
	dbutil "github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
	"github.com/acme/signalforge/internal/shared/validator"
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

func (s *Service) Create(ctx context.Context, policy domain.Policy) (domain.Policy, error) {
	now := s.clock.Now()
	if err := validator.Run(validator.Required("name", policy.Name)); err != nil {
		return domain.Policy{}, err
	}
	if policy.WaitSeconds < 1 {
		policy.WaitSeconds = 300
	}
	if policy.RepeatSeconds < 1 {
		policy.RepeatSeconds = 600
	}
	if policy.MaxRepeats < 0 {
		policy.MaxRepeats = 3
	}
	policy.ID = id.Prefix("esc")
	policy.CreatedAt = now
	policy.UpdatedAt = now
	if err := s.repository.Create(ctx, policy); err != nil {
		return domain.Policy{}, err
	}
	return policy, nil
}

func (s *Service) Update(ctx context.Context, policy domain.Policy) (domain.Policy, error) {
	existing, err := s.repository.FindByID(ctx, policy.ID)
	if err != nil {
		return domain.Policy{}, err
	}
	policy.CreatedAt = existing.CreatedAt
	policy.UpdatedAt = s.clock.Now()
	if err := s.repository.Update(ctx, policy); err != nil {
		return domain.Policy{}, err
	}
	return policy, nil
}

func (s *Service) Get(ctx context.Context, id string) (domain.Policy, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.Policy, int, error) {
	return s.repository.List(ctx, enabledOnly, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Match(ctx context.Context, target matcher.Target) ([]domain.Policy, error) {
	policies, err := s.repository.Active(ctx)
	if err != nil {
		return nil, err
	}
	var matched []domain.Policy
	for _, policy := range policies {
		if policy.Enabled && policy.Matcher.Matches(target) {
			matched = append(matched, policy)
		}
	}
	return matched, nil
}

func (s *Service) CreateMany(ctx context.Context, policies []domain.Policy) ([]domain.Policy, error) {
	created := make([]domain.Policy, 0, len(policies))
	err := dbutil.RunBatch(ctx, policies, func(ctx context.Context, policy domain.Policy) error {
		item, err := s.Create(ctx, policy)
		if err != nil {
			return err
		}
		created = append(created, item)
		return nil
	})
	return created, err
}

func (s *Service) ValidateMany(policies []domain.Policy) error {
	for _, policy := range policies {
		if err := validator.Run(validator.Required("name", policy.Name)); err != nil {
			return err
		}
	}
	return nil
}

func SeverityInRoutes(routes []domain.Route, severity severity.Severity) bool {
	for _, route := range routes {
		if len(route.Severities) == 0 {
			return true
		}
		for _, value := range route.Severities {
			if value == severity.String() {
				return true
			}
		}
	}
	return false
}
