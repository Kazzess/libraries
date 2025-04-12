package redis

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Kazzess/libraries/metrics"
)

var (
	redisAvailability = metrics.NewGaugeVec(
		metrics.GaugeOpts{
			Name: "redis_availability",
			Help: "Indicates the availability of Redis connection (1 for available, 0 for unavailable)",
		},
		[]string{"address", "db"},
	)
)

func checkRedisAvailability(ctx context.Context, client *Client, cfg *Config) {
	redisAvailability.WithLabelValues(cfg.address, strconv.Itoa(cfg.db)).Set(0)

	go func() {
		ticker := time.NewTicker(cfg.health.intervalCheck)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := client.Ping(ctx).Err()
				if err != nil {
					log.Printf("Failed to ping redis server due to error: %v\n", err)

					redisAvailability.WithLabelValues(cfg.address, strconv.Itoa(cfg.db)).Set(0)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, false)
					}
				} else {
					redisAvailability.WithLabelValues(cfg.address, strconv.Itoa(cfg.db)).Set(1)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, true)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
