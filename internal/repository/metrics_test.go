package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/d1zyy/monitor-pc/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
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

func TestMetricsRepository_SaveAndCleanup(t *testing.T) {
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Skip("URL db not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Cleanup(func() {
		_, _ = dbpool.Exec(context.Background(), "DELETE FROM metrics_history")
		dbpool.Close()
	})

	_, err = dbpool.Exec(ctx, "DELETE FROM metrics_history")
	if err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}

	repo := NewMetricsRepository(dbpool)

	expected := &metrics.SystemMetrics{
		CPUPercent: 12.3,
		RAMUsed:    45.2,
		RAMTotal:   32.0,

		DiskUsed:  34.2,
		DiskTotal: 120.0,
	}

	_, err = dbpool.Exec(ctx, `
	INSERT INTO metrics_history
		(cpu_percent, ram_used, ram_total, disk_used, disk_total, created_at)
	VALUES
		($1, $2, $3, $4, $5, NOW() - INTERVAL '40 days')
`,
		10.0,
		1000.0,
		8000.0,
		2000.0,
		100000.0,
	)
	if err != nil {
		t.Fatalf("failed to insert old metric: %v", err)
	}

	err = repo.SaveAndCleanup(ctx, expected)
	if err != nil {
		t.Fatalf("failed to func: %v", err)
	}
	assert.NoError(t, err)

	var count int

	err = dbpool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM metrics_history",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count metrics: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 metrics after cleanup, got: %d", count)
	}

	latest, err := repo.GetLatest(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get last metrics: %v", err)
	}

	if len(latest) != 1 {
		t.Fatalf("expected 1 metrics, got %d", len(latest))
	}

	if latest[0].CPUPercent != expected.CPUPercent ||
		latest[0].RAMUsed != expected.RAMUsed ||
		latest[0].RAMTotal != expected.RAMTotal ||
		latest[0].DiskUsed != expected.DiskUsed ||
		latest[0].DiskTotal != expected.DiskTotal {
		t.Fatalf("retrieved metrics do not match saved metrics")
	}

}
