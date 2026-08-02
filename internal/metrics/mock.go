//test struct for cache

package metrics

type MockCollector struct {
	Count int
}

func (m *MockCollector) Collect() (*SystemMetrics, error) {
	m.Count++

	return &SystemMetrics{
		CPUPercent: float64(m.Count),
		RAMUsed:    float64(m.Count) * 2,
		RAMTotal:   float64(m.Count) * 4,
		DiskUsed:   float64(m.Count) * 3,
		DiskTotal:  float64(m.Count) * 6,
	}, nil
}
