package metrics

type Collector interface {
	Collect() (*SystemMetrics, error)
}
