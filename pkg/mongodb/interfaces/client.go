package interfaces

import (
	"context"

	"github.com/haguru/horus/pkg/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Connect() (*mongo.Client, error)
	CreateSpatialIndex(databaseName string, collectionName string, spatialType string) error
	Delete(databaseName string, collectionName string, id string) error
	Disconnect(context.Context) error
	FindAll(databaseName string, collectionName string) ([]bson.D, error)
	FindOne(databaseName string, collectionName string, id string) bson.D
	InsertRecord(databaseName string, collectionName string, doc interface{}) (string, error)
	Ping(client *mongo.Client) error
	SpaitalQuery(point models.Point, databasName string, collectionName string) ([]bson.D, error)
	Update(databaseName string, collectionName string, id string, crumb interface{}) error
	// SpatialFilter(models.Point) bson.D
}
