package repository

import (
	"context"
	"fmt"

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

func (mr *MetricsRepository) SaveAndCleanup(ctx context.Context, m *metrics.SystemMetrics) error {
	tx, err := mr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	queryInsert := `
		INSERT INTO metrics_history (cpu_percent, ram_used, ram_total, disk_used, disk_total) 
		VALUES ($1, $2, $3, $4, $5)
	`

	queryDelete := `
		DELETE FROM metrics_history WHERE created_at < NOW() - INTERVAL '30 days'
	`
	_, err = tx.Exec(ctx, queryInsert, m.CPUPercent, m.RAMUsed, m.RAMTotal, m.DiskUsed, m.DiskTotal)
	if err != nil {
		return fmt.Errorf("failed to insert metrics: %w", err)
	}

	_, err = tx.Exec(ctx, queryDelete)
	if err != nil {
		return fmt.Errorf("failed to cleanup old metrics: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
