package persistence

import "strings"

func sourceFingerprint(sources []string) string {
	return strings.Join(append([]string(nil), sources...), ",")
}

func mergeSources(left []string, right []string) []string {
	result := append([]string(nil), left...)
	result = append(result, right...)
	return result
}

func sourceDifference(left []string, right []string) []string {
	diff := make([]string, 0, len(left))
	for _, item := range left {
		if !containsSource(right, item) {
			diff = append(diff, item)
		}
	}
	return append([]string(nil), diff...)
}

func containsSource(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
