package http

import (
	"github.com/acme/signalforge/internal/incident/domain"
	"testing"
)

func TestDetachTimelineCopiesMetadata(t *testing.T) {
	events := []domain.TimelineEvent{{Metadata: map[string]any{"owner": "ops"}}}
	copy := detachTimeline(events)
	copy[0].Metadata["owner"] = "changed"
	if events[0].Metadata["owner"] != "ops" {
		t.Fatal("timeline response aliases metadata")
	}
}
