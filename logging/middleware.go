package logging

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Kazzess/libraries/tracing"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	requestIDLogKey = "request_id"
	traceIDLogKey   = "trace_id"
	spanIDLogKey    = "span_id"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mLogger := L(ctx).With(slog.String("endpoint", r.URL.RequestURI()))

		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			mLogger = mLogger.With(slog.String(traceIDLogKey, span.TraceID().String()))
			tracing.TraceValue(ctx, traceIDLogKey, span.TraceID().String())
			mLogger = mLogger.With(slog.String(spanIDLogKey, span.TraceID().String()))
		}

		ctx = ContextWithLogger(ctx, mLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithTraceIDInLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		mLogger := L(ctx).With(slog.String("method", info.FullMethod))

		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			mLogger = mLogger.With(slog.String(traceIDLogKey, span.TraceID().String()))
			tracing.TraceValue(ctx, traceIDLogKey, span.TraceID().String())
			mLogger = mLogger.With(slog.String(spanIDLogKey, span.TraceID().String()))
		}

		ctx = ContextWithLogger(ctx, mLogger)

		return handler(ctx, req)
	}
}
