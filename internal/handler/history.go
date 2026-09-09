package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d1zyy/monitor-pc/internal/metrics"
	"github.com/gin-gonic/gin"
)

type MetricsHistoryHandler struct {
	reader metrics.MetricsHistoryReader
}

func NewMetricsHistoryHandler(reader metrics.MetricsHistoryReader) *MetricsHistoryHandler {
	return &MetricsHistoryHandler{
		reader: reader,
	}
}

func (h *MetricsHistoryHandler) GetLatest(c *gin.Context) {
	limit := 20

	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Limit parameter"})
			return
		}
		limit = parsedLimit
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	history, err := h.reader.GetLatest(ctx, limit)
	if err != nil {
		log.Printf("failed to get latest metrics history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
