package scheduler

type resultCollector struct {
	results map[string]string
}

func newResultCollector() *resultCollector {
	return &resultCollector{results: make(map[string]string)}
}

func (c *resultCollector) Set(name string) {
	c.results[name] = "ok"
}

func (c *resultCollector) Snapshot() map[string]string {
	return c.results
}
