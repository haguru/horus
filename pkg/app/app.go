package app

import (
	"context"
	"fmt"
	"net"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/haguru/horus/config"
	"github.com/haguru/horus/internal/routes"
	pb "github.com/haguru/horus/internal/routes/protos"
	"github.com/haguru/horus/pkg/mongodb"
	"github.com/haguru/horus/pkg/mongodb/interfaces"
	"google.golang.org/grpc"
)

type App struct {
	LoggingClient  logger.LoggingClient
	AppCtx         context.Context
	ServiceConfig  *config.ServiceConfig
	Route          *routes.Route
	GrpcServer     *grpc.Server
	DbServerClient interfaces.Client
}

func NewApp() (*App, error) {
	serviceConfig, err := config.ReadLocalConfig(config.CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read config locally: %v", err)
	}

	lc := logger.NewClient(serviceConfig.Name, serviceConfig.LogLevel)

	db, err := mongodb.NewMongoDB("localhost", 27017, nil)
	if err != nil {
		lc.Errorf("failed to connect, %v\n", err)
		return nil,err
	}

	dbConfig := serviceConfig.Database
	err = db.CreateSpatialIndex(dbConfig.Name, dbConfig.Collection, mongodb.SPATIAL_INDEX_TYPE)
	if err != nil{
		lc.Errorf("failed to create spatial index: %v", err)
		return nil,err
	}
	
	route := routes.NewRoute(lc,&serviceConfig.Database,db)
	return &App{
		LoggingClient:  lc,
		AppCtx:         context.Background(),
		ServiceConfig:  serviceConfig,
		DbServerClient: db,
		Route: route,
	}, nil
}

func (app *App) RunServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.ServiceConfig.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	app.GrpcServer = grpc.NewServer()
	pb.RegisterCrumbDBServer(app.GrpcServer, app.Route)
	app.LoggingClient.Debugf("server listening at %v", lis.Addr())
	err = app.GrpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
