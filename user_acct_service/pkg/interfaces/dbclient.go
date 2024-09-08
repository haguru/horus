package interfaces

import "context"

type DbClient interface {
	Create(databaseName string, collectionName string, doc interface{}) (string, error)
	Delete(databaseName string, collectionName string, filterParms map[string]interface{}) error
	Disconnect(context.Context) error
	DocumentExist(databaseName string, collectionName string, filterParams map[string]interface{}) (bool,error)
	Get(databaseName string, collectionName string, filterParams map[string]interface{}) (interface{}, error)
	Update(databaseName string, collectionName string, filterParams map[string]interface{}, updateType string, items map[string]interface{}) error
}
