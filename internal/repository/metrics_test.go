package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/d1zyy/monitor-pc/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMetricsRepository_SaveAndGetLatest(t *testing.T) {
	connStr := os.Getenv("TEST_DATABASE_URL")

	if connStr == "" {
		t.Skip("TEST_DATABASE_URL is not set. Skip!")
	}

	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := dbPool.Exec(cleanupCtx, "DELETE FROM metrics_history")
		if err != nil {
			t.Logf("failed to clean up metrics_history table: %v", err)
		}

		dbPool.Close()
	})

	_, err = dbPool.Exec(ctx, "DELETE FROM metrics_history")
	if err != nil {
		t.Fatalf("failed to clear metrics_history table: %v", err)
	}

	repo := NewMetricsRepository(dbPool)

	expected := &metrics.SystemMetrics{
		CPUPercent: 50.0,
		RAMUsed:    4096,
		RAMTotal:   8192,
		DiskUsed:   50000,
		DiskTotal:  100000,
	}

	err = repo.Save(ctx, expected)
	if err != nil {
		t.Fatalf("failed to save metrics: %v", err)
	}

	latestMetrics, err := repo.GetLatest(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get latest metrics: %v", err)
	}

	if len(latestMetrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(latestMetrics))
	}

	if latestMetrics[0].CPUPercent != expected.CPUPercent ||
		latestMetrics[0].RAMUsed != expected.RAMUsed ||
		latestMetrics[0].RAMTotal != expected.RAMTotal ||
		latestMetrics[0].DiskUsed != expected.DiskUsed ||
		latestMetrics[0].DiskTotal != expected.DiskTotal ||
		latestMetrics[0].CreatedAt.IsZero() {
		t.Fatalf("retrieved metrics do not match saved metrics")
	}

}
