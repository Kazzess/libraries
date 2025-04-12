package metrics

import (
	"net/http"
	"time"
)

var httpRequestDurationMs = NewHistogramVec(
	HistogramOpts{
		Name:    "http_requests_time_ms",
		Help:    "The duration of HTTP requests in milliseconds",
		Buckets: []float64{10, 50, 100, 500, 1000},
	},
	[]string{"service", "method"},
)

func measureHTTPRequestDurationMs(serviceName, method string, start time.Time) {
	httpRequestDurationMs.
		WithLabelValues(serviceName, method).
		Observe(float64(time.Since(start).Milliseconds()))
}

func RequestDurationMetricHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer measureHTTPRequestDurationMs(r.URL.Path, r.Method, start)
		next.ServeHTTP(w, r)
	})
}

func RequestDurationMetricHTTPMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer measureHTTPRequestDurationMs(r.URL.Path, r.Method, start)
		next.ServeHTTP(w, r)
	}
}
