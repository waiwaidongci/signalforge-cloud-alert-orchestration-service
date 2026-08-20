package domain

import "testing"

func TestTimelineEventCloneCopiesMetadata(t *testing.T) {
	event := TimelineEvent{Metadata: map[string]any{
		"labels": map[string]any{"owner": "ops"},
		"tags":   []string{"urgent", "database"},
	}}
	clone := event.Clone()
	clone.Metadata["labels"].(map[string]any)["owner"] = "changed"
	clone.Metadata["tags"].([]string)[0] = "changed"

	if event.Metadata["labels"].(map[string]any)["owner"] != "ops" {
		t.Fatal("clone exposes nested metadata map")
	}
	if event.Metadata["tags"].([]string)[0] != "urgent" {
		t.Fatal("clone exposes metadata slice")
	}
}
