package repository

import (
	"context"

	"github.com/d1zyy/monitor-pc/internal/metrics"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsRepository struct {
	db *pgxpool.Pool
}

func NewMetricsRepository(db *pgxpool.Pool) *MetricsRepository {
	return &MetricsRepository{db: db}
}

func (mr *MetricsRepository) Save(metrics *metrics.SystemMetrics, ctx context.Context) error {
	_, err := mr.db.Exec(ctx, `
		INSERT INTO metrics_history (cpu_percent, ram_used, ram_total, disk_used, disk_total)
		VALUES ($1, $2, $3, $4, $5)
	`, metrics.CPUPercent, metrics.RAMUsed, metrics.RAMTotal, metrics.DiskUsed, metrics.DiskTotal)

	if err != nil {
		return err
	}
	return nil
}
