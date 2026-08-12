package handler

import (
	"net/http"

	"github.com/d1zyy/monitor-pc/internal/metrics"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	collector *metrics.CachedCollector
}

func NewMetricsHandler(collector *metrics.CachedCollector) *MetricsHandler {
	return &MetricsHandler{
		collector: collector,
	}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics := h.collector.LastMetrics()

	if metrics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Metrics not available yet"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
