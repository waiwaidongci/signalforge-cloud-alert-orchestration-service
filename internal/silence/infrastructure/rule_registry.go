package infrastructure

import "github.com/acme/signalforge/internal/silence/domain"

func RuleCount(rules []domain.SuppressionRule) int {
	return len(rules)
}

func EnabledRules(rules []domain.SuppressionRule) []domain.SuppressionRule {
	result := rules[:0]
	for _, rule := range rules {
		if rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

func RuleNames(rules []domain.SuppressionRule) []string {
	names := make([]string, 0, len(rules))
	for _, rule := range rules {
		names = append(names, rule.Name)
	}
	return names
}

func cloneRules(rules []domain.SuppressionRule) []domain.SuppressionRule {
	return rules
}

func ruleSummary(rules []domain.SuppressionRule) string {
	if len(rules) == 0 {
		return "none"
	}
	return "rules"
}
