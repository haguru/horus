package main

import (
	"context"
	"fmt"

	"github.com/haguru/horus/pkg/app"
)

// type App struct {
// 	// lc            logger.LoggingClient
// 	appCtx        context.Context
// 	serviceConfig config.ServiceConfig
// }

func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Printf("failed to create new app: %v\n", err)
		return
	}
	defer func() {
		if err = app.DbServerClient.Disconnect(context.TODO()); err != nil {
			app.LoggingClient.Errorf("failed to disconnect db server: %v", err)
		}
	}()

	err = app.RunServer()
	if err != nil {
		app.LoggingClient.Errorf("failed to start grpc server: %v", err)
		return
	}
	defer func() {
		app.GrpcServer.Stop()
	}()

	// dbName, collName := "test", "crumbs"
	// err = db.CreateSpatialIndex("test", "crumbs", "2dsphere")
	// if err != nil {
	// 	fmt.Printf("failed to create spatial index, %v\n", err)
	// }

	// data := models.Crumb{
	// 	Location: mongoModels.Point{
	// 		Type:        "Point",
	// 		Coordinates: []float64{-122.64579888741955, 45.691752785517224},
	// 	},
	// 	Message: "testing",
	// 	User:    "bdkmv",
	// }
	// _, err = db.InsertRecord(dbName, collName, data)
	// if err != nil {
	// 	panic(err)
	// }

	// newpoint := mongoModels.Point{
	// 	Type:        "Point",
	// 	Coordinates: []float64{-122.64585552120147, 45.69219926469911},
	// }

	// fmt.Println("spatial query--------------------")

	// dataoutput, err := db.SpaitalQuery(newpoint, dbName, collName)
	// if err != nil {
	// 	panic(err)
	// }

	// d := db.FindOne(dbName, collName, "66cad94c99d236b3235851ab")
	// fmt.Printf("found: %v\n", d)

	// update := models.MessageUpdateRequest{
	// 	Message: "testing_updated",
	// }

	// err = db.Update(dbName, collName, "66cad94c99d236b3235851ab", update)
	// if err != nil {
	// 	panic(err)
	// }

	// err = db.Delete(dbName, collName, "66cadbc22158accbeb4cc24d")
	// if err != nil {
	// 	panic(err)
	// }

	// dataoutput, err = db.SpaitalQuery(newpoint, dbName, collName)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, bsonCrumb := range dataoutput {
	// 	crumb := &models.Crumb{}
	// 	d, err := bson.Marshal(bsonCrumb)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	err = bson.Unmarshal(d, crumb)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(crumb)
	// }

	// fmt.Println("works.....")
}
