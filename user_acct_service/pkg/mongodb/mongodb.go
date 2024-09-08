package mongodb

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/useracctdb/pkg/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Uri        string
	Host       string
	Port       int
	ServerOpts *options.ServerAPIOptions
	Client     *mongo.Client
	lc         logger.LoggingClient
}

const (
	MAXPOOLSIZE = 20
	IDFIELD     = "_id"
)

// NewMongoDB returns a interface for db client and error if it occurs
func NewMongoDB(host string, port int, lc logger.LoggingClient, opts *options.ServerAPIOptions) (interfaces.DbClient, error) {
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

func (db MongoDB) Create(databaseName string, collectionName string, doc interface{}) (string, error) {
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

func (db MongoDB) Get(databaseName string, collectionName string, filterParams map[string]interface{}) (interface{}, error) {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	filter := db.filter(bson.M{}, filterParams)

	results := collection.FindOne(context.TODO(), filter)
	var data bson.D
	err := results.Decode(&data)
	if err != nil {
		db.lc.Errorf("failed to decode results: %v", err)
		return nil, err
	}
	return &data, nil
}

func (db MongoDB) Update(databaseName string, collectionName string, filterParams map[string]interface{}, updateOperator string, items map[string]interface{}) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	filter := db.filter(bson.M{}, filterParams)

	updateItems, err := db.updateCommand(updateOperator,items)
	if err != nil{
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), filter, updateItems)
	if err != nil {
		return err
	}
	return nil
}

func (db MongoDB) Delete(databaseName string, collectionName string, filterParams map[string]interface{}) error {
	collection := db.Client.Database(databaseName).Collection(collectionName)

	filter := db.filter(bson.M{}, filterParams)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	db.lc.Debugf("deleted count: %v\n", res.DeletedCount)
	return nil
}

func (db MongoDB) filter(bsonMap bson.M, searchParams map[string]interface{}) bson.M {
	for key, value := range searchParams {
		bsonMap[key] = value
	}
	return bsonMap
}

func (db MongoDB) DocumentExist(databaseName string, collectionName string, filterParams map[string]interface{}) (bool,error){
	collection := db.Client.Database(databaseName).Collection(collectionName)
	
	filter := db.filter(bson.M{}, filterParams)
	
	found, err := collection.CountDocuments(context.Background(), filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	
	return found > 0, nil
}

func (db MongoDB) updateCommand(updateOperator string, items map[string]interface{}) (interface{},error){
	switch updateOperator {
	case "currentDate":
		return CurrentDateOp{CurrentDate: items}, nil
	case "inc":
		return IncOp{Inc: items}, nil
	case "max":
		return MaxOp{Max: items}, nil
	case "min":
		return MinOp{Min: items}, nil
	case "mul":
		return MulOp{Mul: items}, nil
	case "rename":
		return RenameOp{Rename: items}, nil
	case "set":
		return SetOp{Set: items}, nil
	case "setOnInsert":
		return SetOnInsertOp{SetOnInsert: items}, nil
	case "unset":
		return UnsetOp{Unset: items}, nil
	}

	return nil, fmt.Errorf("failed to create update command")
}
