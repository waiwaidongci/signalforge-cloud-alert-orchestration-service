package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metrics) RecentPathCount(path string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recentPaths[path]
}

func (m *Metrics) RecentPathSnapshot() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	snapshot := make(map[string]int, len(m.recentPaths))
	for path, count := range m.recentPaths {
		snapshot[path] = count
	}
	return snapshot
}
