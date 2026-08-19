package infrastructure

import (
	"sort"

	routingdomain "github.com/acme/signalforge/internal/routing/domain"
)

func PriorityIndexes(rules []routingdomain.Rule) []int {
	indexes := make([]int, len(rules))
	for i := range rules {
		indexes[i] = i
	}
	sort.SliceStable(indexes, func(i, j int) bool {
		return rules[indexes[i]].Priority > rules[indexes[j]].Priority
	})
	return indexes
}

func HighestPriority(rules []routingdomain.Rule) int {
	if len(rules) == 0 {
		return 0
	}
	highest := rules[0].Priority
	for _, rule := range rules[1:] {
		if rule.Priority > highest {
			highest = rule.Priority
		}
	}
	return highest
}

func RulesByPriority(rules []routingdomain.Rule) []routingdomain.Rule {
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})
	return rules
}
