package interfaces

import (
	"context"

	"github.com/haguru/horus/pkg/mogodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Connect() (*mongo.Client, error)
	CreateSpatialIndex(databaseName string, collectionName string, spatialType string) error
	Disconnect(context.Context) error
	InsertRecord(databaseName string, collectionName string, docs []interface{}) error
	Ping(client *mongo.Client) error
	SpaitalQuery(point models.Point, databasName string, collectionName string) ([]bson.D, error)
	FindAll(databaseName string, collectionName string) ([]bson.D, error)
	// SpatialFilter(models.Point) bson.D
}
