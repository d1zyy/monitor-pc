package metrics

import (
	"context"
	"log"
	"sync"
	"time"
)

type CachedCollector struct {
	mu        sync.RWMutex
	last      *SystemMetrics
	collector Collector
	wg        sync.WaitGroup
	//mock      MockCollector
}

func NewCachedCollector(ctx context.Context) (*CachedCollector, error) {
	collector, err := NewCollector()
	if err != nil {
		return nil, err
	}

	c := &CachedCollector{
		collector: collector,
	}

	c.wg.Add(1)
	go c.refreshLoop(ctx)
	return c, nil
}

func (c *CachedCollector) Collect() (*SystemMetrics, error) {
	c.mu.RLock()
	last := c.last
	c.mu.RUnlock()

	if last != nil {
		return last, nil
	}

	metrics, err := c.collector.Collect()
	if err != nil {
		return nil, err
	}

	c.setLast(metrics)
	return metrics, nil
}

func (c *CachedCollector) LastMetrics() *SystemMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.last
}

func (c *CachedCollector) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	defer c.wg.Done()

	for {
		select {
		case <-ticker.C:
			metrics, err := c.collector.Collect()
			if err != nil {
				log.Println("Error collecting metrics:", err)
				continue
			}
			c.setLast(metrics)

		case <-ctx.Done():
			log.Println("Stopping metrics refresh loop")
			return
		}
	}
}

func (c *CachedCollector) setLast(metrics *SystemMetrics) {
	c.mu.Lock()
	c.last = metrics
	c.mu.Unlock()
}

func (c *CachedCollector) Wait() {
	c.wg.Wait()
}
