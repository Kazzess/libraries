package metrics

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var grpcRequestDurationMs = NewHistogramVec(
	HistogramOpts{
		Name:    "grpc_requests_time_ms",
		Help:    "The duration of gRPC requests in milliseconds",
		Buckets: []float64{10, 50, 100, 500, 1000},
	},
	[]string{"service", "method", "is_err"},
)

func measureGRPCRequestDurationMs(serviceName, method string, isErr bool, start time.Time) {
	grpcRequestDurationMs.
		WithLabelValues(serviceName, method, strconv.FormatBool(isErr)).
		Observe(float64(time.Since(start).Milliseconds()))
}

func RequestDurationMetricUnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		defer measureGRPCRequestDurationMs(serviceName, info.FullMethod, err != nil, start)

		return handler(ctx, req)
	}
}
