package infrastructure

import (
	"strings"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
)

func ResolveField(input alertdomain.IngestInput, name string) string {
	switch name {
	case "external_id":
		return input.ExternalID
	case "resource":
		return input.Resource
	case "title":
		return input.Title
	default:
		return ""
	}
}

func DefaultSeverity(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "info"
	}
	return raw
}

func BuildNormalizedLabels(input alertdomain.IngestInput) map[string]string {
	if input.Labels == nil {
		input.Labels = map[string]string{}
	}
	input.Labels["resource"] = ResolveField(input, "resource")
	return input.Labels
}
