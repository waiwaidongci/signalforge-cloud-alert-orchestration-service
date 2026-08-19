package domain

import "testing"

func TestResolveMappingUsesFieldMapping(t *testing.T) {
	source := Source{FieldMapping: map[string]string{"resource": "custom-resource"}}
	if got := ResolveMapping(source, "resource", "default"); got != "custom-resource" {
		t.Fatalf("expected custom-resource, got %q", got)
	}
}
