package infrastructure

import (
	"testing"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	sourcedomain "github.com/acme/signalforge/internal/source/domain"
)

func TestNormalizeInitializesNilMapsBeforeWriting(t *testing.T) {
	input := alertdomain.IngestInput{
		ExternalID: "evt-1",
		Resource:   "api-server-1",
		Severity:   "high",
		Title:      "latency high",
	}
	normalized := Normalize(input, sourcedomain.Source{})
	if normalized.Labels == nil {
		t.Fatal("labels must not be nil after normalization")
	}
	if normalized.Annotations == nil {
		t.Fatal("annotations must not be nil after normalization")
	}
	if normalized.Labels["normalized"] != "true" {
		t.Fatalf("expected normalized marker, got %q", normalized.Labels["normalized"])
	}
}
