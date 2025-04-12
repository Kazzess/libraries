package minio

import (
	"context"
	"time"

	"github.com/Kazzess/libraries/metrics"
)

const minioHealthCheckDuration = 5 * time.Second

var (
	minioAvailability = metrics.NewGaugeVec(
		metrics.GaugeOpts{
			Name: "minio_availability",
			Help: "Indicates the availability of MinIO service (1 for available, 0 for unavailable)",
		},
		[]string{"endpoint"},
	)
)

func checkMinioAvailability(ctx context.Context, client *Client, cfg *Config) error {
	cancelFn, healthErr := client.minioClient.HealthCheck(minioHealthCheckDuration)
	if healthErr != nil {
		return healthErr
	}

	minioAvailability.WithLabelValues(cfg.endpoint).Set(0)

	go func() {
		ticker := time.NewTicker(cfg.health.intervalCheck)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				minioAvailability.WithLabelValues(cfg.endpoint).Set(0)
				cancelFn()

				return
			case <-ticker.C:
				client.minioClient.BucketExists(ctx, "healthcheck")

				if !client.minioClient.IsOffline() {
					minioAvailability.WithLabelValues(cfg.endpoint).Set(1)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, true)
					}
				} else if client.minioClient.IsOffline() {
					minioAvailability.WithLabelValues(cfg.endpoint).Set(0)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, false)
					}
				}
			}
		}
	}()

	return nil
}
