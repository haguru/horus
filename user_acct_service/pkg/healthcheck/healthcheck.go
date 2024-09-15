package healthcheck

import (
	"time"

	"github.com/haguru/horus/useracctdb/config"
	"github.com/haguru/horus/useracctdb/pkg/interfaces"
	appMetrics "github.com/haguru/horus/useracctdb/pkg/prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	PING_INTERVAL = "5s"
	HEALTHY       = 1.0
	UNHEALTHY     = 0.0
)

type HealthCheck struct {
	Health        *health.Server
	ServiceConfig *config.ServiceConfig
	ticker        *time.Ticker
	metrics       *appMetrics.Metrics
}

func NewHealthCheck(config *config.ServiceConfig, metrics *appMetrics.Metrics) (*HealthCheck, error) {
	duration, err := time.ParseDuration(PING_INTERVAL)
	if err != nil {
		return nil, err
	}
	return &HealthCheck{
		Health:        health.NewServer(),
		ServiceConfig: config,
		ticker:        time.NewTicker(duration),
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
	h.Health.SetServingStatus(h.ServiceConfig.Name, status)
}

func (h *HealthCheck) StartHealthCheckService(client interfaces.DbClient) {
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
