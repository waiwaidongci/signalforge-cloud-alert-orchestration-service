package id

import (
	"strings"
	"testing"
)

func TestZeroValueAllocatorUsesFallbackGenerator(t *testing.T) {
	fallback := NewIDGenerator(nil)
	if fallback.Next("probe") == "" {
		t.Fatal("fallback generator must produce an id")
	}
	var allocator Allocator
	value := allocator.AllocateID("evt")
	if !strings.HasPrefix(value, "evt_") || !allocator.Issued(value) {
		t.Fatalf("zero allocator did not issue safe id: %q", value)
	}
	if !AllocationReady(&allocator) {
		t.Fatal("allocator must be ready after first allocation")
	}
	if !AllocationAudited(&allocator, value) {
		t.Fatal("allocated id must be audited")
	}
}
