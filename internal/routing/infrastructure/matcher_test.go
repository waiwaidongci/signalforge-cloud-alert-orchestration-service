package infrastructure

import (
	"testing"

	routingdomain "github.com/acme/signalforge/internal/routing/domain"
)

func TestMatchChannelsDoesNotMutateInput(t *testing.T) {
	input := []routingdomain.Channel{
		{Channel: "email", Destination: "ops@example.com"},
		{Channel: "", Destination: "missing"},
		{Channel: "webhook", Destination: "https://example.com"},
	}
	before := append([]routingdomain.Channel(nil), input...)
	got := MatchChannels(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(got))
	}
	for i := range before {
		if before[i] != input[i] {
			t.Fatalf("input was mutated at %d: before=%+v after=%+v", i, before[i], input[i])
		}
	}
}
