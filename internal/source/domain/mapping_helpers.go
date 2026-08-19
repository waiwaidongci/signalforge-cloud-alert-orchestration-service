package domain

func ResolveMapping(source Source, target string, fallback string) string {
	if value := source.FieldMapping[target]; value != "" {
		return value
	}
	return fallback
}
