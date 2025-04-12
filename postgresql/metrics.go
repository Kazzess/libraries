package psql

import (
	"context"
	"log"
	"net/url"
	"time"

	"git.adapticode.com/libraries/golang/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	postgresAvailability = metrics.NewGaugeVec(
		metrics.GaugeOpts{
			Name: "postgres_availability",
			Help: "Indicates the availability of PostgreSQL connection (1 for available, 0 for unavailable)",
		},
		[]string{"host", "database"},
	)
)

func checkPostgresAvailability(ctx context.Context, pool *pgxpool.Pool, cfg *Config) error {
	u, err := url.Parse(cfg.dsn)
	if err != nil {
		return err
	}

	host := u.Hostname()

	pathSegments := u.Path
	if len(pathSegments) > 1 {
		pathSegments = pathSegments[1:]
	}

	postgresAvailability.WithLabelValues(host, pathSegments).Set(0)

	go func() {
		ticker := time.NewTicker(cfg.health.intervalCheck)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pingErr := pool.Ping(ctx)
				if pingErr != nil {
					log.Printf("Failed to ping PostgreSQL due to error: %v\n", pingErr)

					postgresAvailability.WithLabelValues(host, pathSegments).Set(0)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, false)
					}
				} else {
					postgresAvailability.WithLabelValues(host, pathSegments).Set(1)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, true)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
