package application

import (
	"context"
	"time"

	"github.com/acme/signalforge/internal/dedup/domain"
)

type Resolver interface {
	Resolve(ctx context.Context, input domain.Input) (string, error)
}

type Service struct {
	window time.Duration
}

func NewService(window time.Duration) *Service {
	if window <= 0 {
		window = 15 * time.Minute
	}
	return &Service{window: window}
}

func (s *Service) Resolve(_ context.Context, input domain.Input) (string, error) {
	return domain.Fingerprint(input), nil
}

func (s *Service) Window() time.Duration {
	return s.window
}
