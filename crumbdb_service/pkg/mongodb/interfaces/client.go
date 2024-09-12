package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	// Connect returns a mongodb client and error.
	// If an error occurs mongodb client will be nil
	Connect() (*mongo.Client, error)

	// CreateSpatialIndex returns error if client is unable to create a spatial index
	// this is needed to search database by (longitude, latitude) coordinates
	CreateSpatialIndex(databaseName string, collectionName string, spatialType string) error

	// Delete removes a document from the database. Returns nil error if successful
	Delete(databaseName string, collectionName string, id string) error

	// Disconnect returns error if client is unable to disconnect from mongodb
	Disconnect(context.Context) error

	// FindAll retrieves all documents in the database. Returns an array of bson.D and error.
	// if an error occurs then a nil is return and an error
	FindAll(databaseName string, collectionName string) ([]bson.D, error)

	// FindOne retrieves a document by ID. Returns a bson.D
	FindOne(databaseName string, collectionName string, id string) (*bson.D, error)

	// InsertRecord returns ID, as string, and error.
	// if error occurs an empty string is returned along with the error
	InsertRecord(databaseName string, collectionName string, doc interface{}) (string, error)

	// Ping returns error if mongodb is unreachable
	Ping(client *mongo.Client) error

	// SpaitalQuery queries database for data based on coordinates. Returns array of bson.D and error
	// if error occurs a nil is returned as well as an error
	SpaitalQuery(pointType string, coordinates []float64, databaseName string, collectionName string) ([]bson.D, error)

	// Update modifies a document given a ID. Returns a nil error when sucessful
	Update(databaseName string, collectionName string, id string, items map[string]interface{}) error
}
