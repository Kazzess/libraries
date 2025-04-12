package metrics

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		host:              "localhost",
		port:              8080,
		readTimeout:       10 * time.Second,
		writeTimeout:      10 * time.Second,
		readHeaderTimeout: 5 * time.Second,
	}
	server, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	if server == nil {
		t.Fatal("Expected non-nil server instance")
	}
}

func TestServerRunAndClose(t *testing.T) {
	cfg := Config{
		host:              "localhost",
		port:              8080,
		readTimeout:       10 * time.Second,
		writeTimeout:      10 * time.Second,
		readHeaderTimeout: 5 * time.Second,
	}
	server, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := server.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Failed to run server: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(1 * time.Second)

	// Now close the server
	if err := server.Close(); err != nil {
		t.Errorf("Failed to close server: %v", err)
	}

	cancel()
}

func TestHTTPServerResponse(t *testing.T) {
	cfg := Config{
		host:              "localhost",
		port:              8080,
		readTimeout:       10 * time.Second,
		writeTimeout:      10 * time.Second,
		readHeaderTimeout: 5 * time.Second,
	}

	server, err := NewServer(&cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err2 := server.Run(ctx)
		if err2 != nil {
			t.Errorf("Failed to run server: %v", err2)
		}
	}()

	time.Sleep(1 * time.Second)

	// Use httptest to create a request to the server
	req, _ := http.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	// Directly use the server's handler to serve the request
	server.httpServer.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	all, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	strings.Contains(string(all), "go_gc_duration_seconds")
	strings.Contains(string(all), "go_goroutines")
	strings.Contains(string(all), "go_memstats_alloc_bytes")
	strings.Contains(string(all), "promhttp_metric_handler_requests_in_flight")
	strings.Contains(string(all), "promhttp_metric_handler_requests_total")
}
