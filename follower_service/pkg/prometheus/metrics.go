package prometheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/followerdb/config"
	"go.opentelemetry.io/otel/trace"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const (
	REQUEST_COUNTER      = "request_counter"
	ERROR_COUNTER        = "error_counter"
	HELP_REQUEST_COUNTER = "Number of requests received"
	HELP_ERROR_COUNTER   = "Number of errors that occured"
	METRICSPORT             = 52112
	
)

var(
	BUCKETS = []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}
)

type Metrics struct {
	Registry      *prometheus.Registry
	// RequestsCount prometheus.Counter
	// ErrorCount    prometheus.Counter
	GrpcMetrics   *grpc_prometheus.ServerMetrics
	MetricServer  *http.Server
}

func NewMetrics(config *config.ServiceConfig) *Metrics {
	serverMetrics := grpc_prometheus.NewServerMetrics(
		grpc_prometheus.WithServerHandlingTimeHistogram(
			grpc_prometheus.WithHistogramBuckets(BUCKETS),
		),
	)
	metrics := &Metrics{
		Registry: prometheus.NewRegistry(),
		// RequestsCount: prometheus.NewCounter(prometheus.CounterOpts{
		// 	Namespace: config.Name,
		// 	Name:      REQUEST_COUNTER,
		// 	Help:      HELP_REQUEST_COUNTER,
		// }),
		// ErrorCount: prometheus.NewCounter(prometheus.CounterOpts{
		// 	Namespace: config.Name,
		// 	Name:      ERROR_COUNTER,
		// 	Help:      HELP_ERROR_COUNTER,
		// }),
		GrpcMetrics: serverMetrics,
	}
	metrics.Registry.MustRegister(metrics.GrpcMetrics)
	prometheusHandler := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})
	metrics.MetricServer = &http.Server{Handler: prometheusHandler, Addr: fmt.Sprintf("0.0.0.0:%d", METRICSPORT)}

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

// func Recovery(p any) (err error) {
// 	panicsTotal.Inc()
// 	level.Error(rpcLogger).Log("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
// 	return status.Errorf(codes.Internal, "%s", p)
// }
