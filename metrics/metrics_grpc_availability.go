package metrics

import (
	"context"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	gRPCConnectionAvailability = NewGaugeVec(
		GaugeOpts{
			Name: "grpc_connection_availability",
			Help: "Current availability of the service (1 = available, 0 = unavailable)",
		},
		[]string{"from_service", "to_service"},
	)
)

type GRPCService interface {
	Connection() grpc.ClientConnInterface
}

type GRPCConnectionMonitor struct {
	service     GRPCService
	pingTimer   time.Duration
	serviceFrom string
	serviceTo   string

	available int32
	cancel    context.CancelFunc
}

func NewGRPCConnectionMonitor(
	service GRPCService,
	pingTimer time.Duration,
	serviceFrom string,
	serviceTo string,
) *GRPCConnectionMonitor {
	return &GRPCConnectionMonitor{
		service:     service,
		pingTimer:   pingTimer,
		serviceFrom: serviceFrom,
		serviceTo:   serviceTo,
	}
}

func (s *GRPCConnectionMonitor) Start(ctx context.Context) {
	healthClient := grpc_health_v1.NewHealthClient(s.service.Connection())

	gRPCConnectionAvailability.WithLabelValues(s.serviceFrom, s.serviceTo).Set(0)

	checkConnection := func() {
		var status grpc_health_v1.HealthCheckResponse_ServingStatus

		check, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: "",
		})
		if err != nil {
			status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
		} else {
			status = check.GetStatus()
		}

		switch status {
		case grpc_health_v1.HealthCheckResponse_SERVING:
			if atomic.CompareAndSwapInt32(&s.available, 0, 1) {
				gRPCConnectionAvailability.WithLabelValues(s.serviceFrom, s.serviceTo).Set(1)
			}
		default:
			if atomic.CompareAndSwapInt32(&s.available, 1, 0) {
				gRPCConnectionAvailability.WithLabelValues(s.serviceFrom, s.serviceTo).Set(0)
			}
		}
	}

	checkConnection()

	ticker := time.NewTicker(s.pingTimer)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				checkConnection()
			}
		}
	}()

	tStop := func() {
		ticker.Stop()
	}

	s.cancel = tStop
}

func (s *GRPCConnectionMonitor) Close() error {
	s.cancel()

	return nil
}
