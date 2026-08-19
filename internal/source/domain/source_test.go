package domain

import "testing"

func TestSourceMappingHandlesNilFieldMapping(t *testing.T) {
	source := Source{}
	if got := source.Mapping("resource"); got != "resource" {
		t.Fatalf("expected resource, got %q", got)
	}
}
