package application

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/silence/domain"
)

func missingSilence() *domain.Silence {
	return nil
}

func hasSilenceMatch(silences []domain.Silence, target matcher.Target, at time.Time) bool {
	return ActiveSilence(silences, target, at) != nil
}
