package infrastructure

import (
	"context"
	"github.com/acme/signalforge/internal/incident/domain"
	"testing"
)

func TestTimelineListDoesNotExposeInternalSlice(t *testing.T) {
	s := NewInMemoryTimeline()
	event := domain.TimelineEvent{IncidentID: "i", Message: "x", Metadata: map[string]any{"k": "v"}}
	_ = s.Append(context.Background(), event)
	event.Metadata["k"] = "changed after append"
	p, _, _ := s.List(context.Background(), "i", 10, 0)
	p[0].Message = "changed"
	p[0].Metadata["k"] = "changed"
	q, _, _ := s.List(context.Background(), "i", 10, 0)
	if q[0].Message != "x" || q[0].Metadata["k"] != "v" {
		t.Fatalf("mutated: %#v", q)
	}
}

func TestTimelineListRejectsInvalidWindow(t *testing.T) {
	s := NewInMemoryTimeline()
	_ = s.Append(context.Background(), domain.TimelineEvent{IncidentID: "i"})
	if _, _, err := s.List(context.Background(), "i", -1, 0); err == nil {
		t.Fatal("negative limit must be rejected")
	}
}
