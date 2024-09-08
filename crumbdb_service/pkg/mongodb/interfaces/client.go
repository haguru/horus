package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Connect() (*mongo.Client, error)
	CreateSpatialIndex(databaseName string, collectionName string, spatialType string) error
	Delete(databaseName string, collectionName string, id string) error
	Disconnect(context.Context) error
	FindAll(databaseName string, collectionName string) ([]bson.D, error)
	FindOne(databaseName string, collectionName string, id string) (*bson.D, error)
	InsertRecord(databaseName string, collectionName string, doc interface{}) (string, error)
	Ping(client *mongo.Client) error
	SpaitalQuery(coordinates []float64, databasName string, collectionName string) ([]bson.D, error)
	Update(databaseName string, collectionName string, id string, items map[string]interface{}) error
	// SpatialFilter(models.Point) bson.D
}
