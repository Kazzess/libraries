package healthcheck

import (
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const defaultCheckInterval = 10 * time.Second

type GRPCHealthServer struct {
	*health.Server
	healthState struct {
		mu        *sync.RWMutex
		statusMap map[string]bool
		val       bool
	}
}

func NewGRPCHealthServer() *GRPCHealthServer {
	return &GRPCHealthServer{
		Server: health.NewServer(),
		healthState: struct {
			mu        *sync.RWMutex
			statusMap map[string]bool
			val       bool
		}{
			mu:        &sync.RWMutex{},
			statusMap: make(map[string]bool),
			val:       false,
		},
	}
}

func (gs *GRPCHealthServer) HealthCheck(ctx context.Context, checkInterval time.Duration) {
	if checkInterval <= 0 {
		checkInterval = defaultCheckInterval
	}

	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				gs.healthState.mu.RLock()
				gs.healthState.val = true

				for _, status := range gs.healthState.statusMap {
					if !status {
						gs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
						gs.healthState.val = false
						break
					}
				}

				if gs.healthState.val {
					gs.SetServingStatus(
						"", grpc_health_v1.HealthCheckResponse_SERVING)
				}

				gs.healthState.mu.RUnlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (gs *GRPCHealthServer) SetStatus(dependencyName string, status bool) {
	gs.healthState.mu.Lock()
	defer gs.healthState.mu.Unlock()

	gs.healthState.statusMap[dependencyName] = status
}
