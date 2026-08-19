package scheduler

import "sync"

type resultCollector struct {
	mu      sync.Mutex
	results map[string]string
}

func newResultCollector() *resultCollector {
	return &resultCollector{results: make(map[string]string)}
}

func (c *resultCollector) Set(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results[name] = "ok"
}

func (c *resultCollector) Snapshot() map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.results
}
