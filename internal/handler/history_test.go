package handler

import (
	"context"
	"time"

	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/d1zyy/monitor-pc/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockHistoryReader struct {
	metrics []*metrics.HistoryMetrics
	err     error
	limit   int
	called  bool
}

func (m *mockHistoryReader) GetLatest(ctx context.Context, limit int) ([]*metrics.HistoryMetrics, error) {
	m.limit = limit
	m.called = true

	if m.err != nil {
		return nil, m.err
	}

	return m.metrics, nil
}

func TestMetricsHistoryReader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := &mockHistoryReader{
		metrics: []*metrics.HistoryMetrics{
			{
				SystemMetrics: metrics.SystemMetrics{
					CPUPercent: 10.5,
					RAMUsed:    2048,
					RAMTotal:   8192,
					DiskUsed:   50000,
					DiskTotal:  100000,
				},
				CreatedAt: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			},
			{
				SystemMetrics: metrics.SystemMetrics{
					CPUPercent: 20.0,
					RAMUsed:    4096,
					RAMTotal:   8192,
					DiskUsed:   60000,
					DiskTotal:  100000,
				}, CreatedAt: time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC),
			},
		},
		err:   nil,
		limit: 2,
	}

	h := NewMetricsHistoryHandler(reader)

	router := gin.New()
	router.GET("/history", h.GetLatest)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/history?limit=2", nil)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, 2, reader.limit)
	assert.JSONEq(t, `[
		{
			"cpu_percent": 10.5,
			"ram_used": 2048,
			"ram_total": 8192,
			"disk_used": 50000,
			"disk_total": 100000,
			"created_at": "2024-06-01T10:00:00Z"
		},
		{
			"cpu_percent": 20.0,
			"ram_used": 4096,
			"ram_total": 8192,
			"disk_used": 60000,
			"disk_total": 100000,
			"created_at": "2024-06-01T11:00:00Z"
		}
	]`, rec.Body.String())
}

func TestMetricsHistoryHandler_RepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := &mockHistoryReader{
		err: errors.New("database error"),
	}

	h := NewMetricsHistoryHandler(reader)

	router := gin.New()
	router.GET("/history", h.GetLatest)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/history", nil)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"error": "Failed to retrieve metrics history"}`, rec.Body.String())
}

func TestMetricsHistoryHandler_InvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := &mockHistoryReader{}

	h := NewMetricsHistoryHandler(reader)

	router := gin.New()
	router.GET("/history", h.GetLatest)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/history?limit=invalid", nil)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, reader.called)
	assert.JSONEq(t, `{"error": "Invalid Limit parameter"}`, rec.Body.String())
}
