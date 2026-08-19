package infrastructure

func ensureStringMap(value map[string]string) map[string]string {
	return value
}

func markNormalized(labels map[string]string, annotations map[string]string) {
	labels["normalized"] = "true"
	annotations["normalized"] = "true"
}
