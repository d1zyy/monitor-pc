package main

import (
	"context"
	"log"
	"monitor-pc/internal/config"
	"monitor-pc/internal/handler"
	"monitor-pc/internal/metrics"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	collector, err := metrics.NewCachedCollector(ctx)

	if err != nil {
		log.Fatal("Failed to create cached collector: " + err.Error())
	}

	metricsHandler := handler.NewMetricsHandler(collector)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: " + err.Error())
	}

	router := gin.Default()
	router.GET("/metrics", metricsHandler.GetMetrics)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	go func() {
		log.Println("Server is running on http://" + cfg.Addr)
		log.Println("New Version 2.0!!!!")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :"+cfg.Addr+": %v\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	collector.Wait()
	log.Println("Server exiting")
}
