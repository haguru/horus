package healthcheck

import (
	"time"

	"github.com/haguru/horus/crumbdb/config"
	"github.com/haguru/horus/crumbdb/pkg/mongodb/interfaces"
	appMetrics "github.com/haguru/horus/crumbdb/pkg/prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	HEALTHY   = 1.0
	UNHEALTHY = 0.0
)

type HealthCheck struct {
	Health        *health.Server
	ServiceConfig *config.ServiceConfig
	ticker        *time.Ticker
	metrics       *appMetrics.Metrics
}

func NewHealthCheck(config *config.ServiceConfig, metrics *appMetrics.Metrics, pingInterval time.Duration) (*HealthCheck, error) {
	return &HealthCheck{
		Health:        health.NewServer(),
		ServiceConfig: config,
		ticker:        time.NewTicker(pingInterval),
		metrics:       metrics,
	}, nil
}

func (h *HealthCheck) Initialize(serviceGrpcServer *grpc.Server) {
	// Register the health service with the gRPC server
	healthpb.RegisterHealthServer(serviceGrpcServer, h.Health)

	// Set the service health status
	h.SetStatus(healthpb.HealthCheckResponse_SERVING)
}

func (h *HealthCheck) SetStatus(status healthpb.HealthCheckResponse_ServingStatus) {
	h.Health.SetServingStatus(h.ServiceConfig.ServiceName, status)
}

func (h *HealthCheck) StartHealthCheckService(client interfaces.Client) {
	for {
		select {
		case <-h.ticker.C:
			err := client.Ping()
			if err != nil {
				h.metrics.HealthMetric.Set(UNHEALTHY)
				h.SetStatus(healthpb.HealthCheckResponse_NOT_SERVING)
			} else {
				h.metrics.HealthMetric.Set(HEALTHY)
				h.SetStatus(healthpb.HealthCheckResponse_SERVING)
			}
		}
	}
}
