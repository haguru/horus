package mogodb

import (
	"context"
	"fmt"

	"github.com/haguru/horus/pkg/mogodb/interfaces"
	"github.com/haguru/horus/pkg/mogodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MAX_DISTANCE      = 100
	SPATIAL_INDEX_KEY = "location"
)

type MongoDB struct {
	Uri        string
	Host       string
	Port       int
	ServerOpts *options.ServerAPIOptions
	Client     *mongo.Client
	// databaseName string
	// context    context.Context
}

func NewMongoDB(host string, port int, opts *options.ServerAPIOptions) (interfaces.Client, error) {
	db := &MongoDB{
		Host:       host,
		Port:       port,
		ServerOpts: opts,
	}
	client, err := db.Connect()
	if err != nil {
		return nil, err
	}
	db.Client = client

	return db, nil
}

func (db MongoDB) Connect() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1).SetStrict(true).SetDeprecationErrors(true)
	if db.ServerOpts != nil {
		serverAPI = db.ServerOpts
	}
	uri := fmt.Sprintf("mongodb://%v:%v/?maxPoolSize=20&w=majority",db.Host,db.Port)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Creat new client
	var err error
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = db.Ping(client)
	if err != nil {
		return nil, fmt.Errorf("failed to successfully ping mongodb server: %v", err)
	}

	return client, nil
}

func (db MongoDB) Ping(client *mongo.Client) error {
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (db MongoDB) Disconnect(context context.Context) error {
	if err := db.Client.Disconnect(context); err != nil {
		return err
	}

	return nil
}

func (db MongoDB) CreateSpatialIndex(databaseName string, collectionName string, spatialType string) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: SPATIAL_INDEX_KEY, Value: spatialType}},
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (db MongoDB) InsertRecord(databaseName string, collectionName string, docs []interface{}) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	_, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}

	return nil
}

func (db MongoDB) SpaitalQuery(point models.Point, databasName string, collectionName string) ([]bson.D, error) {
	filter := db.spatialFilter(point)
	collection := db.Client.Database(databasName).Collection(collectionName)

	var docs []bson.D

	output, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	err = output.All(context.TODO(), &docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (db MongoDB) FindAll(databaseName string, collectionName string) ([]bson.D, error) {
	collection := db.Client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	var results []bson.D
	for cur.Next(context.TODO()) {
		// Create a value into which the single document can be decoded
		var elem bson.D
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}

func (db MongoDB) spatialFilter(point models.Point) bson.D {
	return bson.D{
		{Key: SPATIAL_INDEX_KEY, Value: bson.D{
			{Key: "$near", Value: bson.D{
				{Key: "$geometry", Value: point},
				{Key: "$maxDistance", Value: MAX_DISTANCE},
			}},
		}},
	}
}
