package application

import (
	"sync"
	"testing"

	"github.com/acme/signalforge/internal/dedup/domain"
)

func TestIndexOwnsEntriesAndSnapshots(t *testing.T) {
	index := NewIndex()
	original := domain.IndexEntry{Fingerprint: "a", Labels: map[string]string{"zone": "east"}}
	index.Register("a", original)
	original.Labels["zone"] = "west"

	entry, ok := index.Lookup("a")
	if !ok || entry.Labels["zone"] != "east" {
		t.Fatalf("registered entry must not share caller labels: %#v", entry)
	}
	entry.Labels["zone"] = "north"
	again, _ := index.Lookup("a")
	if again.Labels["zone"] != "east" {
		t.Fatalf("lookup must not expose stored labels: %#v", again)
	}
	snapshot := index.Snapshot()
	snapshot["a"] = domain.IndexEntry{Fingerprint: "changed", Labels: map[string]string{"zone": "south"}}
	snapshot["a"].Labels["zone"] = "south"
	final, _ := index.Lookup("a")
	if final.Fingerprint != "a" || final.Labels["zone"] != "east" {
		t.Fatalf("snapshot must be isolated from index: %#v", final)
	}
}

func TestIndexConcurrentSnapshotAndRegister(t *testing.T) {
	index := NewIndex()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			index.Register("shared", domain.IndexEntry{Fingerprint: "shared", Labels: map[string]string{"zone": "east"}})
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			for range index.Snapshot() {
			}
		}
	}()
	close(start)
	wg.Wait()
}
