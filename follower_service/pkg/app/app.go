package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	// "time"

	"github.com/haguru/horus/followerdb/config"
	"github.com/haguru/horus/followerdb/internal/routes"
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"
	"github.com/haguru/horus/followerdb/pkg/interfaces"
	"github.com/haguru/horus/followerdb/pkg/mongodb"
	appMetrics "github.com/haguru/horus/followerdb/pkg/prometheus"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/health"
	// "google.golang.org/grpc/health/grpc_health_v1"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	// healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	METRICS_ENDPOINT = "/metrics"
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

	// Create a gRPC Server with gRPC interceptor.
	app.GrpcServer = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(app.metrics.GrpcMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(appMetrics.ExemplarFromContext))),
			logging.UnaryServerInterceptor(appMetrics.InterceptorLogger(app.LoggingClient), logging.WithFieldsFromContext(appMetrics.LogTraceID)),
		),
		grpc.ChainStreamInterceptor(
			grpc.StreamServerInterceptor(app.metrics.GrpcMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(appMetrics.ExemplarFromContext))),
			logging.StreamServerInterceptor(appMetrics.InterceptorLogger(app.LoggingClient), logging.WithFieldsFromContext(appMetrics.LogTraceID)),
		),
	)

	pb.RegisterFollowerDBServer(app.GrpcServer, app.Route)
	app.metrics.GrpcMetrics.InitializeMetrics(app.GrpcServer)

	// // health service
	// healthService := health.NewServer()
	// // Register the health service with the gRPC server
	// grpc_health_v1.RegisterHealthServer(app.GrpcServer, healthService)
	// // Set the service health status
	// healthService.SetServingStatus(app.ServiceConfig.Name, grpc_health_v1.HealthCheckResponse_SERVING)

	// ticker := time.NewTicker(5 * time.Second)

	// go func(app *App) {
	// 	for {
	// 		select {
	// 		// case <-done:
	// 		//     return
	// 		case <-ticker.C:
	// 			err := app.DbServerClient.Ping()
	// 			if err != nil {
	// 				// app.LoggingClient.Debug("updating health status: NOT_SERVING",)
	// 				healthService.SetServingStatus(app.ServiceConfig.Name, healthpb.HealthCheckResponse_NOT_SERVING)
	// 			}else{
	// 				// app.LoggingClient.Debug("updating health status: SERVING")
	// 				healthService.SetServingStatus(app.ServiceConfig.Name,healthpb.HealthCheckResponse_SERVING)
	// 			}

				
	// 		}
	// 	}
	// }(app)

	go func() {
		metricsServer := &http.Server{Addr: fmt.Sprintf(":%d", app.ServiceConfig.Metrics.Port)}
		muxHandler := http.NewServeMux()
		muxHandler.Handle(METRICS_ENDPOINT, promhttp.HandlerFor(app.metrics.Registry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}))

		metricsServer.Handler = muxHandler
		app.LoggingClient.Debugf("server(prometheus) listening at %v", app.ServiceConfig.Metrics.Port)
		if err := metricsServer.ListenAndServe(); err != nil {
			app.LoggingClient.Error("failed to start prometheus client")
		}
	}()

	app.LoggingClient.Debugf("server(GRPC) listening at %v", lis.Addr())
	err = app.GrpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
