package main

import (
	"fmt"

	"github.com/haguru/horus/gateway/pkg/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Printf("failed to create new app: %v\n", err)
		return
	}
	// defer func() {
	// 	if err = app.DbServerClient.Disconnect(context.TODO()); err != nil {
	// 		app.LoggingClient.Errorf("failed to disconnect db server: %v", err)
	// 	}
	// }()

	err = app.RunServer()
	if err != nil {
		app.LoggingClient.Errorf("failed to start app: %v", err)
		return
	}
	defer func() {
		app.GrpcServer.Stop()
	}()
}
