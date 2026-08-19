package application

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/silence/domain"
)

func selectSilence(silences []domain.Silence, target matcher.Target, at time.Time) *domain.Silence {
	for i := range silences {
		if silences[i].Active(at) && silences[i].Matcher.Matches(target) {
			return &silences[i]
		}
	}
	return nil
}

func suppressionCandidates(rules []domain.SuppressionRule, target matcher.Target) []domain.SuppressionRule {
	result := make([]domain.SuppressionRule, 0, len(rules))
	for _, rule := range rules {
		if rule.Enabled && rule.SourceMatcher.Matches(target) {
			result = append(result, rule)
		}
	}
	return result
}
