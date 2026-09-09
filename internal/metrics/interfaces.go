package metrics

import "context"

type MetricsSaver interface {
	Save(ctx context.Context, metrics *SystemMetrics) error
}

type MetricsHistoryReader interface {
	GetLatest(ctx context.Context, limit int) ([]*HistoryMetrics, error)
}
