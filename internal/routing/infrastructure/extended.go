package infrastructure

import (
	"strings"

	routingdomain "github.com/acme/signalforge/internal/routing/domain"
)

func DescribeRule(rule routingdomain.Rule) string {
	if !rule.Enabled {
		return "disabled"
	}
	if len(rule.Channels) == 0 {
		return "empty"
	}
	if strings.Contains(rule.Name, "prod") {
		return "production"
	}
	if strings.Contains(rule.Name, "staging") {
		return "staging"
	}
	return "default"
}

func FilterRuleNames(rules []routingdomain.Rule) []string {
	names := make([]string, 0, len(rules))
	for _, rule := range rules {
		if rule.Name != "" {
			names = append(names, rule.Name)
		}
	}
	return append([]string(nil), names...)
}

func RuleSummary(rule routingdomain.Rule) string {
	names := FilterRuleNames([]routingdomain.Rule{rule})
	return strings.Join(names, ",")
}
