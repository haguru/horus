package prometheus

import (
	"context"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/followerdb/config"
	"go.opentelemetry.io/otel/trace"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var BUCKETS = []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}

type Metrics struct {
	Registry     *prometheus.Registry
	HealthMetric prometheus.Gauge
	GrpcMetrics  *grpc_prometheus.ServerMetrics
}

func NewMetrics(config *config.ServiceConfig) *Metrics {
	healthMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: config.Name,
			Name:      "health",
			Help:      "Checks the health of the connection to DB",
		})
	serverMetrics := grpc_prometheus.NewServerMetrics(
		grpc_prometheus.WithServerHandlingTimeHistogram(
			grpc_prometheus.WithHistogramBuckets(BUCKETS),
		),
	)

	metrics := &Metrics{
		Registry:     prometheus.NewRegistry(),
		HealthMetric: healthMetric,
		GrpcMetrics:  serverMetrics,
	}

	metrics.Registry.MustRegister(metrics.GrpcMetrics, healthMetric)

	return metrics
}

func Auth(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	// TODO: This is example only, perform proper Oauth/OIDC verification!
	if token != "yolo" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}
	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
	return ctx, nil
}

func Health(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
}

func InterceptorLogger(lc logger.LoggingClient) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			lc.Debug(msg)
		case logging.LevelInfo:
			lc.Info(msg)
		case logging.LevelWarn:
			lc.Warn(msg)
		case logging.LevelError:
			lc.Error(msg)
		default:
			lc.Debugf("unknown level %v", lvl)
		}
	})
}

func LogTraceID(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func ExemplarFromContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"traceID": span.TraceID().String()}
	}
	return nil
}
