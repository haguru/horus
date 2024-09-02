package mongodb

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"go.mongodb.org/mongo-driver/bson"
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
)

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

func (db MongoDB) Create() error {
	return nil
}

func (db MongoDB) GetUser() error {
	return nil
}

func (db MongoDB) UpdatePassword() error {
	return nil
}

func (db MongoDB) DeleteUser() error {
	return nil
}
