package mongodb

import (
	"context"
	"fmt"

	"github.com/haguru/horus/crumbdb/pkg/mongodb/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
)

const (
	MAX_DISTANCE       = 100
	MAXPOOLSIZE        = 20
	SPATIAL_INDEX_TYPE = "2dsphere"
	SPATIAL_INDEX_KEY  = "location"
	_ID                = "_id"
)

type MongoDB struct {
	Uri        string
	Host       string
	Port       int
	ServerOpts *options.ServerAPIOptions
	Client     *mongo.Client
	lc         logger.LoggingClient
}

// NewMongoDB returns a interface for db client and error if it occurs
func NewMongoDB(host string, port int, lc logger.LoggingClient, opts *options.ServerAPIOptions) (interfaces.Client, error) {
	db := &MongoDB{
		Host:       host,
		Port:       port,
		lc:         lc,
		ServerOpts: opts,
	}
	client, err := db.Connect()
	if err != nil {
		return nil, err
	}
	db.Client = client

	return db, nil
}

// Connect returns a mongodb client and error.
// If an error occurs mongodb client will be nil
func (db MongoDB) Connect() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1).SetStrict(true).SetDeprecationErrors(true)
	if db.ServerOpts != nil {
		serverAPI = db.ServerOpts
	}
	uri := fmt.Sprintf("mongodb://%v:%v/?maxPoolSize=%v&w=majority", db.Host, db.Port, MAXPOOLSIZE)
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

// Ping returns error if mongodb is unreachable
func (db MongoDB) Ping(client *mongo.Client) error {
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return err
	}

	return nil
}

// Disconnect returns error if client is unable to disconnect from mongodb
func (db MongoDB) Disconnect(context context.Context) error {
	if err := db.Client.Disconnect(context); err != nil {
		return err
	}

	return nil
}

// CreateSpatialIndex returns error if client is unable to create a spatial index
// this is needed to search database by (longitude, latitude) coordinates
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

// InsertRecord returns ID, as string, and error.
// if error occurs an empty string is returned along with the error
func (db MongoDB) InsertRecord(databaseName string, collectionName string, doc interface{}) (string, error) {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	r, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return "", err
	}

	objId, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to get objectID")
	}

	return objId.String(), nil
}

// SpaitalQuery queries database for data based on coordinates. Returns array of bson.D and error
// if error occurs a nil is returned as well as an error
func (db MongoDB) SpaitalQuery(point interface{}, databasName string, collectionName string) ([]bson.D, error) {
	filter := db.spatialFilter(point)
	collection := db.Client.Database(databasName).Collection(collectionName)

	output, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var docs []bson.D
	err = output.All(context.TODO(), &docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

// FindAll retrieves all documents in the database. Returns an array of bson.D and error.
// if an error occurs then a nil is return and an error
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

// FindOne retrieves a document by ID. Returns a bson.D
func (db MongoDB) FindOne(databaseName string, collectionName string, id string) (*bson.D, error) {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	// get bson id filter
	objid := db.idFilter(id)

	results := collection.FindOne(context.TODO(), objid)
	var data bson.D
	err := results.Decode(&data)
	if err != nil {
		db.lc.Errorf("failed to decode results: %v", err)
		return nil, err
	}
	return &data, nil
}

// Update modifies a document given a ID. Returns a nil error when sucessful
func (db MongoDB) Update(databaseName string, collectionName string, id string, crumb interface{}) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)
	c, err := db.setMessageUpdate(crumb)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(context.TODO(), objectID, c)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a document from the database. Returns nil error if successful
func (db MongoDB) Delete(databaseName string, collectionName string, id string) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)
	res, err := collection.DeleteOne(context.TODO(), db.idFilter(id))
	if err != nil {
		return err
	}

	db.lc.Debugf("deleted count: %v\n", res.DeletedCount)
	return nil
}

func (db MongoDB) spatialFilter(point interface{}) bson.D {
	return bson.D{
		{Key: SPATIAL_INDEX_KEY, Value: bson.D{
			{Key: "$near", Value: bson.D{
				{Key: "$geometry", Value: point},
				{Key: "$maxDistance", Value: MAX_DISTANCE},
			}},
		}},
	}
}

func (db MongoDB) idFilter(id string) bson.D {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	return bson.D{{Key: _ID, Value: objectID}}
}

func (db MongoDB) setMessageUpdate(data interface{}) (bson.D, error) {
	bsonData := bson.D{}
	err := db.convertToBson(data, &bsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data to bson.D")
	}

	return bson.D{
		{Key: "$set", Value: bsonData},
	}, nil
}

func (db MongoDB) convertToBson(value interface{}, doc *bson.D) error {
	data, err := bson.Marshal(value)
	if err != nil {
		return err
	}

	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return err
	}

	return nil
}
