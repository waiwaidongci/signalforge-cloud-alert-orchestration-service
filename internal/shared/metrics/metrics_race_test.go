package metrics

import (
	"sync"
	"testing"
)

func TestObserveHTTPConcurrentRecentPathAccounting(t *testing.T) {
	m := New()
	const workers = 8
	const calls = 100
	var wg sync.WaitGroup
	ready := make(chan struct{}, workers)
	start := make(chan struct{})
	stopReads := make(chan struct{})
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		<-start
		for {
			select {
			case <-stopReads:
				return
			default:
				_ = m.RecentPathSnapshot()
			}
		}
	}()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ready <- struct{}{}
			<-start
			for j := 0; j < calls; j++ {
				m.ObserveHTTP("GET", "/alerts", "200", 0.001)
			}
		}()
	}
	for i := 0; i < workers; i++ {
		<-ready
	}
	close(start)
	wg.Wait()
	close(stopReads)
	<-readDone
	if got := m.RecentPathCount("/alerts"); got != workers*calls {
		t.Fatalf("recent path count = %d, want %d", got, workers*calls)
	}
	if got := m.RecentPathSnapshot()["/alerts"]; got != workers*calls {
		t.Fatalf("snapshot count = %d", got)
	}
}
