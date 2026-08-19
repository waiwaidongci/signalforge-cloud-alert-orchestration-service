package infrastructure

import "github.com/acme/signalforge/internal/silence/domain"

func RuleSummary(rules []domain.SuppressionRule) string {
	if len(rules) == 0 {
		return "none"
	}
	if len(rules) < 3 {
		return "few"
	}
	return "many"
}

func EnabledRuleCount(rules []domain.SuppressionRule) int {
	count := 0
	for _, rule := range rules {
		if rule.Enabled {
			count++
		}
	}
	return count
}

func RuleLabels(rules []domain.SuppressionRule) []string {
	labels := make([]string, 0, len(rules))
	for _, rule := range rules {
		labels = append(labels, rule.Name)
	}
	return append([]string(nil), labels...)
}
