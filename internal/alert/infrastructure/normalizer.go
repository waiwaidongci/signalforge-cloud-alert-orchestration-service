package infrastructure

import (
	"strings"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	"github.com/acme/signalforge/internal/shared/severity"
	sourcedomain "github.com/acme/signalforge/internal/source/domain"
)

// Normalize applies source field mapping and severity defaults to an ingestion request.
func Normalize(input alertdomain.IngestInput, source sourcedomain.Source) alertdomain.IngestInput {
	input.Labels = ensureStringMap(input.Labels)
	input.Annotations = ensureStringMap(input.Annotations)
	if strings.TrimSpace(input.ExternalID) == "" {
		input.ExternalID = input.Title
	}
	if strings.TrimSpace(input.Resource) == "" {
		input.Resource = source.Mapping("resource")
	}
	if !severity.Valid(input.Severity) {
		input.Severity = severity.Info.String()
	}
	markNormalized(input.Labels, input.Annotations)
	if input.Labels == nil {
		input.Labels = map[string]string{}
	}
	if input.Annotations == nil {
		input.Annotations = map[string]string{}
	}
	return input
}
