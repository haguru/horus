package interfaces

import "context"

type DbClient interface {
	// Connect returns a mongodb client and error.
	// If an error occurs mongodb client will be nil
	Create(databaseName string, collectionName string, doc interface{}) (string, error)
	
	// Ping returns error if mongodb is unreachable
	Delete(databaseName string, collectionName string, filterParms map[string]interface{}) error

	// Disconnect returns error if client is unable to disconnect from mongodb
	Disconnect(context.Context) error

	// DocumentExist checks to see if a document exists in database. Returns bool and error if client fails to run command.
	DocumentExist(databaseName string, collectionName string, filterParams map[string]interface{}) (bool,error)

	// Get reteives a document from database. Returns an interface containing the document and error if client fails to decode data. 
	Get(databaseName string, collectionName string, filterParams map[string]interface{}) (interface{}, error)

	// Update updates a single document in database. Returns error if client fails to  update document or build update command
	Update(databaseName string, collectionName string, filterParams map[string]interface{}, updateType string, items map[string]interface{}) error
}
