package infrastructure

import (
	"testing"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
)

func TestBuildNormalizedLabelsInitializesMap(t *testing.T) {
	input := alertdomain.IngestInput{Resource: "api-server-1"}
	got := BuildNormalizedLabels(input)
	if got == nil {
		t.Fatal("expected normalized labels to be initialized")
	}
	if got["resource"] != "api-server-1" {
		t.Fatalf("expected resource mapping, got %q", got["resource"])
	}
}
