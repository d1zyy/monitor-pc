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

func (mr *MetricsRepository) Save(ctx context.Context, metrics *metrics.SystemMetrics) error {
	_, err := mr.db.Exec(ctx, `
		INSERT INTO metrics_history (cpu_percent, ram_used, ram_total, disk_used, disk_total)
		VALUES ($1, $2, $3, $4, $5)
	`, metrics.CPUPercent, metrics.RAMUsed, metrics.RAMTotal, metrics.DiskUsed, metrics.DiskTotal)

	if err != nil {
		return err
	}
	return nil
}

func (mr *MetricsRepository) GetLatest(ctx context.Context, limit int) ([]*metrics.HistoryMetrics, error) {
	rows, err := mr.db.Query(ctx, `
		SELECT cpu_percent, ram_used, ram_total, disk_used, disk_total, created_at
		FROM metrics_history
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*metrics.HistoryMetrics
	for rows.Next() {
		var m metrics.HistoryMetrics
		if err := rows.Scan(&m.CPUPercent, &m.RAMUsed, &m.RAMTotal, &m.DiskUsed, &m.DiskTotal, &m.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
