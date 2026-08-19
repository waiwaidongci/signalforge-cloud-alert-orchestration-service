package domain

func ResolveMapping(source Source, target string, fallback string) string {
	if source.FieldMapping == nil {
		source.FieldMapping = map[string]string{}
	}
	if value := source.FieldMapping[target]; value != "" {
		return value
	}
	return fallback
}
