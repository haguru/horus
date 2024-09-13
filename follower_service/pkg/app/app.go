package app

import (
	"context"
	"fmt"
	"net"

	"github.com/haguru/horus/followerdb/config"
	"github.com/haguru/horus/followerdb/internal/routes"
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"
	"github.com/haguru/horus/followerdb/pkg/interfaces"
	"github.com/haguru/horus/followerdb/pkg/mongodb"
	appMetrics "github.com/haguru/horus/followerdb/pkg/prometheus"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	promPort = 52112
)

type App struct {
	AppCtx         context.Context
	DbServerClient interfaces.DbClient
	GrpcServer     *grpc.Server
	LoggingClient  logger.LoggingClient
	Route          *routes.Route
	ServiceConfig  *config.ServiceConfig
	metrics        *appMetrics.Metrics
}

func NewApp() (*App, error) {
	serviceConfig, err := config.ReadLocalConfig(config.CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read config locally: %v", err)
	}

	validate := validator.New()
	err = validate.Struct(serviceConfig)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	lc := logger.NewClient(serviceConfig.Name, serviceConfig.LogLevel)

	host := serviceConfig.Database.Host
	port := serviceConfig.Database.Port
	db, err := mongodb.NewMongoDB(host, port, lc, nil)
	if err != nil {
		lc.Errorf("failed to connect, %v\n", err)
		return nil, err
	}

	metrics := appMetrics.NewMetrics(serviceConfig)

	// initiate routes
	route := routes.NewRoute(lc, &serviceConfig.Database, db, validate, metrics)

	// prometheus example
	// http.Handle("/metrics", prometheusHandler)

	return &App{
		LoggingClient:  lc,
		AppCtx:         context.Background(),
		ServiceConfig:  serviceConfig,
		DbServerClient: db,
		Route:          route,
		metrics:        metrics,
	}, nil
}

func (app *App) RunServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.ServiceConfig.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// TODO move to metrics
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	// TODO move to metrics
	// Create a gRPC Server with gRPC interceptor.
	app.GrpcServer = grpc.NewServer(

		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			grpc.UnaryServerInterceptor(app.metrics.GrpcMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext))),
			logging.UnaryServerInterceptor(appMetrics.InterceptorLogger(app.LoggingClient), logging.WithFieldsFromContext(appMetrics.LogTraceID)),
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(appMetrics.Auth), selector.MatchFunc(appMetrics.Health)),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			grpc.StreamServerInterceptor(app.metrics.GrpcMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext))),
			logging.StreamServerInterceptor(appMetrics.InterceptorLogger(app.LoggingClient), logging.WithFieldsFromContext(appMetrics.LogTraceID)),
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(appMetrics.Auth), selector.MatchFunc(appMetrics.Health)),
			// recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	// app.GrpcServer = grpc.NewServer()
	pb.RegisterFollowerDBServer(app.GrpcServer, app.Route)
	app.metrics.GrpcMetrics.InitializeMetrics(app.GrpcServer)

	// TODO move to metrics
	go func(app *App) {
		app.LoggingClient.Debugf("server(prometheus) listening at %v", app.metrics.MetricServer.Addr)
		if err := app.metrics.MetricServer.ListenAndServe(); err != nil {
			app.LoggingClient.Error("failed to start prometheus client")
		}
	}(app)

	app.LoggingClient.Debugf("server(GRPC) listening at %v", lis.Addr())
	err = app.GrpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
