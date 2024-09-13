package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/haguru/horus/useracctdb/config"
	pb "github.com/haguru/horus/useracctdb/internal/routes/protos"
	"github.com/haguru/horus/useracctdb/pkg/interfaces"
)

const (
	MIN_USERNAME_LEN = 1
	MIN_PASSWORD_LEN = 10
	UPDATE_OPERATOR  = "set"
)

type Route struct {
	dbConfig  *config.Database
	dbClient  interfaces.DbClient
	lc        logger.LoggingClient
	validator *validator.Validate

	pb.UnimplementedUserAcctDBServer
}

func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.DbClient, validator *validator.Validate) *Route {
	return &Route{
		dbConfig:  config,
		dbClient:  dbclient,
		lc:        lc,
		validator: validator,
	}
}

func (r *Route) Create(ctx context.Context, user *pb.User) (*pb.Id, error) {
	// Validate the User struct
	err := r.validator.Struct(user)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	// Verify user does not exist
	filterParams := map[string]interface{}{"email": user.GetEmail()}
	exist, err := r.dbClient.DocumentExist(r.dbConfig.Name, r.dbConfig.Collection, filterParams)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("user with email address exists")
	}

	id := &pb.Id{}
	id.Value, err = r.dbClient.Create(r.dbConfig.Name, r.dbConfig.Collection, user)
	if err != nil {
		return nil, fmt.Errorf("database failed to create user: %v", err)
	}

	return id, nil
}

func (r *Route) GetUser(ctx context.Context, userReq *pb.UserRequest) (*pb.User, error) {
	// Validate the UserRequest struct
	err := r.validator.Struct(userReq)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		return nil, fmt.Errorf("validation error: %s", errors)
	}

	user := &pb.User{}
	filterParams := map[string]interface{}{"email": userReq.GetEmail()}
	res, err := r.dbClient.Get(r.dbConfig.Name, r.dbConfig.Collection, filterParams)
	if err != nil {
		return nil, fmt.Errorf("database failed to retrieve user data: %v", err)
	}

	err = r.toUser(user, res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return user, nil
}

func (r *Route) UpdatePassword(ctx context.Context, passwdReq *pb.PasswordRequest) (*pb.Status, error) {
	status := &pb.Status{}
	// Validate the UserRequest struct
	err := r.validator.Struct(passwdReq)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		status.Value = http.StatusBadRequest
		return status, fmt.Errorf("validation error: %s", errors)
	}

	filterParams := map[string]interface{}{"email": passwdReq.Email}
	updateItem := map[string]interface{}{"password": passwdReq.Password}
	err = r.dbClient.Update(r.dbConfig.Name, r.dbConfig.Collection, filterParams, UPDATE_OPERATOR, updateItem)
	if err != nil {
		status.Value = http.StatusInternalServerError
		return status, fmt.Errorf("database failed to update password: %v", err)
	}
	status.Value = http.StatusOK

	return status, nil
}

func (r *Route) Delete(ctx context.Context, userReq *pb.UserRequest) (*pb.Status, error) {
	status := &pb.Status{}
	// Validate the UserRequest struct
	err := r.validator.Struct(userReq)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		status.Value = http.StatusBadRequest
		return status, fmt.Errorf("validation error: %s", errors)
	}
	filterParams := map[string]interface{}{"email": userReq.GetEmail()}
	err = r.dbClient.Delete(r.dbConfig.Name, r.dbConfig.Collection, filterParams)
	if err != nil {
		status.Value = http.StatusInternalServerError
		return status, fmt.Errorf("database failed to delete user: %v", err)
	}

	status.Value = http.StatusOK

	return status, nil
}

func (r *Route) toUser(user *pb.User, doc interface{}) error {
	data, err := bson.Marshal(doc)
	if err != nil {
		return err
	}

	err = bson.Unmarshal(data, user)
	if err != nil {
		return err
	}

	return nil
}
