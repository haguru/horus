package interfaces

import "context"

type DbClient interface {
	Create(databaseName string, collectionName string, doc interface{}) (string, error)
	Delete(databaseName string, collectionName string, filterParms map[string]interface{}) error
	Disconnect(context.Context) error
	Get(databaseName string, collectionName string, filterParams map[string]interface{}) (interface{}, error)
	Update(databaseName string, collectionName string, filterParms map[string]interface{}, updateItems map[string]interface{}) error
}
