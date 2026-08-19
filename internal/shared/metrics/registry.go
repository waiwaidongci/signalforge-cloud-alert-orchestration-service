package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}
