package domain

import "testing"

func TestRuleCloneChannelsIsIndependent(t *testing.T) {
	rule := Rule{Channels: []Channel{{Channel: "email", Destination: "ops@example.com"}}}
	cloned := rule.CloneChannels()
	cloned[0].Destination = "changed@example.com"
	if rule.Channels[0].Destination != "ops@example.com" {
		t.Fatalf("original channels were mutated: %+v", rule.Channels)
	}
}
