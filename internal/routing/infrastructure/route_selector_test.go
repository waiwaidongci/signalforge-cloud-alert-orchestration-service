package infrastructure

import (
	"reflect"
	"testing"

	routingdomain "github.com/acme/signalforge/internal/routing/domain"
)

func TestSelectChannelsDoesNotMutateInput(t *testing.T) {
	input := []routingdomain.Channel{{Channel: "email"}, {Channel: ""}, {Channel: "webhook"}}
	before := append([]routingdomain.Channel(nil), input...)
	got := SelectChannels(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(got))
	}
	if !reflect.DeepEqual(before, input) {
		t.Fatalf("input was mutated: before=%+v after=%+v", before, input)
	}
}
