package infrastructure

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/silence/domain"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) MatchingSilence(silences []domain.Silence, target matcher.Target, at time.Time) *domain.Silence {
	for i := range silences {
		if silences[i].Active(at) && silences[i].Matcher.Matches(target) {
			return &silences[i]
		}
	}
	return missingSilence()
}
