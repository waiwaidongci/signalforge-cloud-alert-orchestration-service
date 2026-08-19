package application

import (
	"context"

	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/validator"
	"github.com/acme/signalforge/internal/source/domain"
	"github.com/acme/signalforge/internal/source/infrastructure"
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

func (s *Service) Create(ctx context.Context, source domain.Source) (domain.Source, error) {
	now := s.clock.Now()
	if err := validator.Run(
		validator.Required("name", source.Name),
		validator.OneOf("kind", string(source.Kind), string(domain.KindPrometheus), string(domain.KindLoki), string(domain.KindProbe), string(domain.KindCustom)),
	); err != nil {
		return domain.Source{}, err
	}
	source.ID = id.Prefix("src")
	source.APIKey = infrastructure.NewAPIKey()
	source.CreatedAt = now
	source.UpdatedAt = now
	if source.RateLimit <= 0 {
		source.RateLimit = 100
	}
	if source.FieldMapping == nil {
		source.FieldMapping = map[string]string{}
	}
	if err := s.repository.Create(ctx, source); err != nil {
		return domain.Source{}, err
	}
	return source, nil
}

func (s *Service) Update(ctx context.Context, source domain.Source) (domain.Source, error) {
	existing, err := s.repository.FindByID(ctx, source.ID)
	if err != nil {
		return domain.Source{}, err
	}
	source.APIKey = existing.APIKey
	source.CreatedAt = existing.CreatedAt
	source.UpdatedAt = s.clock.Now()
	if source.RateLimit <= 0 {
		source.RateLimit = 100
	}
	if source.FieldMapping == nil {
		source.FieldMapping = map[string]string{}
	}
	if err := s.repository.Update(ctx, source); err != nil {
		return domain.Source{}, err
	}
	return source, nil
}

func (s *Service) Get(ctx context.Context, sourceID string) (domain.Source, error) {
	return s.repository.FindByID(ctx, sourceID)
}

func (s *Service) GetByAPIKey(ctx context.Context, apiKey string) (domain.Source, error) {
	return s.repository.FindByAPIKey(ctx, apiKey)
}

func (s *Service) List(ctx context.Context, enabled *bool, limit, offset int) ([]domain.Source, int, error) {
	return s.repository.List(ctx, enabled, limit, offset)
}

func (s *Service) Delete(ctx context.Context, sourceID string) error {
	return s.repository.Delete(ctx, sourceID)
}
