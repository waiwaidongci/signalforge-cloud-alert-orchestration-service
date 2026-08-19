package infrastructure

import "strings"

func ruleFingerprint(rules []string) string {
	return strings.Join(append([]string(nil), rules...), ",")
}

func mergeRules(left []string, right []string) []string {
	result := append([]string(nil), left...)
	result = append(result, right...)
	return result
}

func ruleDifference(left []string, right []string) []string {
	diff := make([]string, 0, len(left))
	for _, item := range left {
		if !containsRule(right, item) {
			diff = append(diff, item)
		}
	}
	return append([]string(nil), diff...)
}

func containsRule(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
