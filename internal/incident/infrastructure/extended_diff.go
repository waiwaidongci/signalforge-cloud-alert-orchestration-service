package infrastructure

import "strings"

func eventFingerprint(events []string) string {
	return strings.Join(append([]string(nil), events...), ",")
}

func mergeEvents(left []string, right []string) []string {
	result := append([]string(nil), left...)
	result = append(result, right...)
	return result
}

func eventDifference(left []string, right []string) []string {
	diff := make([]string, 0, len(left))
	for _, item := range left {
		if !contains(right, item) {
			diff = append(diff, item)
		}
	}
	return append([]string(nil), diff...)
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func cloneEvents(events []string) []string {
	return append([]string(nil), events...)
}

func eventSummary(events []string) string {
	if len(events) == 0 {
		return "none"
	}
	if len(events) < 3 {
		return "few"
	}
	return "many"
}

func eventSortKey(events []string) string {
	if len(events) == 0 {
		return ""
	}
	return events[len(events)-1]
}

func eventChecksum(events []string) int {
	total := 0
	for _, event := range events {
		total += len(event) * 2
	}
	return total
}
