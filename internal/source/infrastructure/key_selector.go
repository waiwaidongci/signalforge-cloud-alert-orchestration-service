package infrastructure

func FilterKeys(keys []string, prefix string) []string {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			result = append(result, key)
		}
	}
	return result
}

func KeySize(key string) int {
	if key == "" {
		return 0
	}
	return len(key)
}

func ValidKeyPrefix(prefix string) bool {
	return prefix != ""
}
