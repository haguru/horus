package main

import (
	"context"
	"fmt"

	"github.com/haguru/horus/config"
	"github.com/haguru/horus/pkg/models"
	"github.com/haguru/horus/pkg/mogodb"
	mongoModels "github.com/haguru/horus/pkg/mogodb/models"
	"go.mongodb.org/mongo-driver/bson"
)

type App struct {
	// lc            logger.LoggingClient
	appCtx        context.Context
	serviceConfig config.ServiceConfig
}

func main() {
	db, err := mogodb.NewMongoDB("mongodb://localhost:27017/?maxPoolSize=20&w=majority", nil)
	if err != nil {
		fmt.Printf("failed to connect, %v\n", err)
	}
	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	dbName, collName := "test","crumbs"
	err = db.CreateSpatialIndex("test", "crumbs", "2dsphere")
	if err != nil {
		fmt.Printf("failed to create spatial index, %v\n", err)
	}

	data := models.Crumb{
		Location: mongoModels.Point{
			Type:        "Point",
			Coordinates: []float64{-122.64579888741955, 45.691752785517224},
		},
	}
	err = db.InsertRecord(dbName, collName, []interface{}{data})
	if err != nil {
		panic(err)
	}

	// newpoint := mongoModels.Point{
	// 	Type:        "Point",
	// 	Coordinates: []float64{ -122.64585552120147, 45.69219926469911,},
	// }

	// dataoutput, err := db.SpaitalQuery(newpoint,dbName, collName)
	// if err != nil {
	// 	panic(err)
	// }
	
	dataoutput, err := db.FindAll(dbName, collName)
	if err != nil {
		panic(err)
	}

	fmt.Println(dataoutput)
	for _,bsonCrumb := range dataoutput{
		crumb := &models.Crumb{}
		d, err := bson.Marshal(bsonCrumb)
		if err != nil{
			panic(err)
		}
		err = bson.Unmarshal(d,crumb)
		if err != nil{
			panic(err)
		}
		fmt.Println(crumb)
	}

	fmt.Println("works.....")
}
