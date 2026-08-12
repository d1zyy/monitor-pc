package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d1zyy/monitor-pc/internal/config"
	"github.com/d1zyy/monitor-pc/internal/handler"
	"github.com/d1zyy/monitor-pc/internal/metrics"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	collector, err := metrics.NewCachedCollector(ctx)

	if err != nil {
		log.Fatal("Failed to create cached collector: " + err.Error())
	}

	metricsHandler := handler.NewMetricsHandler(collector)
	healthHandler := handler.NewHealthHandler()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: " + err.Error())
	}

	// Set Gin
	router := gin.Default()
	router.GET("/metrics", metricsHandler.GetMetrics)
	router.GET("/health", healthHandler.GetHealth)
	router.GET("/version", handler.GetVersion)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("Server is running on http://" + cfg.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("could not listen on %s: %w", cfg.Addr, err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal")
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
		stop()
	}

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("Server close error: %v", closeErr)
		}
	}

	collector.Wait()
	log.Println("Server exiting")
}
