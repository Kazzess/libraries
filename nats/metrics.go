package mynats

import (
	"context"
	"strconv"
	"strings"
	"time"

	"git.adapticode.com/libraries/golang/metrics"
	"github.com/nats-io/nats.go"
)

type ObserveWithErr func(err *error)

var (
	deliveryBuckets       = []float64{1, 10, 25, 50, 100, 150, 200, 500, 1000, 2500, 5000, 10000, 20000, 30000}
	processMessageBuckets = []float64{1, 10, 25, 50, 100, 150, 200, 500, 1000, 2500, 5000, 10000, 20000, 30000}
)

var (
	// deliveryTimeMs is a histogram that measures the time for message delivery from NATS (milliseconds).
	deliveryTimeMs = metrics.NewHistogramVec(
		metrics.HistogramOpts{
			Name:    "nats_stream_delivery_time_ms",
			Help:    "The time for message delivery from NATS (milliseconds)",
			Buckets: deliveryBuckets,
		},
		[]string{
			"subject",
			"consumer_id",
			"is_async",
		},
	)

	// messageProcessingTimeMs is a histogram that measures the time that consumer process message from consume
	//till ack (milliseconds).
	messageProcessingTimeMs = metrics.NewHistogramVec(
		metrics.HistogramOpts{
			Name:    "nats_stream_message_processing_time_ms",
			Help:    "The time that consumer process message from consume till ack (milliseconds)",
			Buckets: processMessageBuckets,
		},
		[]string{
			"subject",
			"consumer_id",
			"is_async",
			"is_err",
		},
	)

	// natsAvailability is a gauge that indicates the availability of NATS connection
	//(1 for connected, 0 for disconnected).
	natsAvailability = metrics.NewGaugeVec(
		metrics.GaugeOpts{
			Name: "nats_availability",
			Help: "Indicates the availability of NATS connection (1 for connected, 0 for disconnected)",
		},
		[]string{"endpoint", "consumer_id"},
	)
)

// ObserveDeliveryTimeMs observes the time for message delivery from NATS.
func ObserveDeliveryTimeMs(subject, consumerID string, ts time.Time, isAsync bool) {
	deliveryTimeMs.
		WithLabelValues(subject, consumerID, strconv.FormatBool(isAsync)).
		Observe(float64(time.Since(ts).Nanoseconds()))
}

// ObserveProcessingTimeMs observes the time that consumer process message from consume till ack.
func ObserveProcessingTimeMs(subject, consumerID string, isAsync bool) ObserveWithErr {
	ts := time.Now()

	return func(err *error) {
		var isErr bool
		if *err != nil {
			isErr = true
		}

		observer := messageProcessingTimeMs.WithLabelValues(
			subject,
			consumerID,
			strconv.FormatBool(isAsync),
			strconv.FormatBool(isErr),
		)
		observer.Observe(float64(time.Since(ts).Nanoseconds()))
	}
}

// metricsMetadata returns the timestamp and consumer ID from the message metadata.
func metricsMetadata(msg *Msg) (time.Time, string) {
	metadata, err := msg.Metadata()
	if err != nil {
		return time.Time{}, ""
	}

	return metadata.Timestamp, metadata.Consumer
}

// MonitorNatsAvailability monitors the availability of NATS connection.
func checkNatsAvailability(ctx context.Context, client *nats.Conn, cfg *Config) {
	dsn := strings.Join(cfg.servers, ",")

	natsAvailability.WithLabelValues(dsn, cfg.consumerID).Set(0)

	go func() {
		ticker := time.NewTicker(cfg.health.intervalCheck)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if client.IsConnected() {
					natsAvailability.WithLabelValues(dsn, cfg.consumerID).Set(1)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, true)
					}
				} else {
					natsAvailability.WithLabelValues(dsn, cfg.consumerID).Set(0)

					if cfg.health.checker != nil {
						cfg.health.checker.SetStatus(cfg.health.name, false)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return
}
