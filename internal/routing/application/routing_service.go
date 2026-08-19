package application

import (
	"context"
	"sort"

	"github.com/acme/signalforge/internal/routing/domain"
	routinginfra "github.com/acme/signalforge/internal/routing/infrastructure"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/matcher"
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

func (s *Service) Create(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	now := s.clock.Now()
	if err := validator.Run(validator.Required("name", rule.Name)); err != nil {
		return domain.Rule{}, err
	}
	rule.ID = id.Prefix("rte")
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.repository.Create(ctx, rule); err != nil {
		return domain.Rule{}, err
	}
	return rule, nil
}

func (s *Service) Update(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	existing, err := s.repository.FindByID(ctx, rule.ID)
	if err != nil {
		return domain.Rule{}, err
	}
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = s.clock.Now()
	if err := s.repository.Update(ctx, rule); err != nil {
		return domain.Rule{}, err
	}
	return rule, nil
}

func (s *Service) Get(ctx context.Context, id string) (domain.Rule, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.Rule, int, error) {
	return s.repository.List(ctx, enabledOnly, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Resolve(ctx context.Context, target matcher.Target) ([]domain.Channel, error) {
	rules, err := s.repository.Active(ctx)
	if err != nil {
		return nil, err
	}
	rules = append([]domain.Rule(nil), rules...)
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].Priority > rules[j].Priority })
	for _, rule := range rules {
		if rule.Enabled && rule.Match.Matches(target) {
			return routinginfra.MatchChannels(rule.CloneChannels()), nil
		}
	}
	return nil, nil
}
