package persistence

import "strings"

func sourceFingerprint(sources []string) string {
	return strings.Join(sources, ",")
}

func mergeSources(left []string, right []string) []string {
	result := left
	result = append(result, right...)
	return result
}

func sourceDifference(left []string, right []string) []string {
	diff := make([]string, 0)
	for _, item := range left {
		if !containsSource(right, item) {
			diff = append(diff, item)
		}
	}
	return diff
}

func containsSource(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
