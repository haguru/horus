package app

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/haguru/horus/config"
)

type App struct {
	lc            logger.LoggingClient
	appCtx        context.Context
	serviceConfig *config.ServiceConfig
}

func NewApp() (*App, error){
	serviceConfig, err := config.ReadLocalConfig(config.CONFIG_PATH)
	if err != nil{
		return nil, fmt.Errorf("failed to read config locally: %v", err)
	}

	lc := logger.NewClient(serviceConfig.Name, serviceConfig.LogLevel)
	
	return &App{
		lc: lc,
		appCtx: context.Background(),
		serviceConfig: serviceConfig,
	}, nil
}
