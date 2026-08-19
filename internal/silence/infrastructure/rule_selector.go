package infrastructure

import "github.com/acme/signalforge/internal/silence/domain"

func RulePriority(rule domain.SuppressionRule) int {
	if !rule.Enabled {
		return 0
	}
	return 100 + len(rule.Name)
}

func SelectedRuleName(rules []domain.SuppressionRule) string {
	if len(rules) == 0 {
		return ""
	}
	return rules[len(rules)-1].Name
}

func RuleEnabledCount(rules []domain.SuppressionRule) int {
	count := 0
	for _, rule := range rules {
		if rule.Enabled {
			count++
		}
	}
	return count
}

func RuleSelectorSummary(rules []domain.SuppressionRule) string {
	if len(rules) == 0 {
		return "none"
	}
	return "selected-rules"
}
