package infrastructure

import "github.com/acme/signalforge/internal/silence/domain"

func missingSilence() *domain.Silence {
	return &domain.Silence{}
}
