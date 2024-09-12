package routes

import (
	"context"
	"fmt"

	"github.com/haguru/horus/crumbdb/config"
	pb "github.com/haguru/horus/crumbdb/internal/routes/protos"
	"github.com/haguru/horus/crumbdb/pkg/mongodb/interfaces"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type Route struct {
	dbConfig  *config.Database
	dbClient  interfaces.Client
	lc        logger.LoggingClient
	validator *validator.Validate
	pb.UnimplementedCrumbDBServer
}

func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.Client) *Route {
	return &Route{
		dbConfig:  config,
		dbClient:  dbclient,
		lc:        lc,
		validator: validator.New(),
	}
}

func (r *Route) Create(ctx context.Context, crumb *pb.Crumb) (*pb.Id, error) {
	r.lc.Debugf("received Create request: %v", crumb)

	// Validate the User struct
	err := r.validator.Struct(crumb)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	id, err := r.dbClient.InsertRecord(r.dbConfig.Name, r.dbConfig.Collection, crumb)
	if err != nil {
		return nil, err
	}

	return &pb.Id{Value: id}, nil
}

func (r *Route) GetCrumbs(point *pb.Point, stream pb.CrumbDB_GetCrumbsServer) error {
	r.lc.Debug("received new GetCumbs request")

	// Validate the User struct
	err := r.validator.Struct(point)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return fmt.Errorf("validation error: %s", errors)
	}

	data, err := r.dbClient.SpaitalQuery(point.Type, point.GetCoordinates(), r.dbConfig.Name, r.dbConfig.Collection)
	if err != nil {
		r.lc.Errorf("failed to run spatial query: %v", err)
		return err
	}
	for _, item := range data {
		// unmarshall data to grpc data type
		doc, err := bson.Marshal(item)
		if err != nil {
			r.lc.Errorf("failed to marshal an item in data: %v", err)
			return err
		}

		crumb := &pb.Crumb{}
		err = bson.Unmarshal(doc, crumb)
		if err != nil {
			r.lc.Errorf("failed to unmarshal an item in data: %v", err)
			return err
		}

		// send crumb
		err = stream.Send(crumb)
		if err != nil {
			r.lc.Errorf("failed to send item in data: %v", err)
			return err
		}

	}
	return nil
}

func (r *Route) Update(ctx context.Context, crumb *pb.Crumb) (*pb.Id, error) {
	r.lc.Debug("received new update request")

	messageItem := map[string]interface{}{"message": crumb.Message}

	err := r.dbClient.Update(r.dbConfig.Name, r.dbConfig.Collection, crumb.Id, messageItem)
	if err != nil {
		r.lc.Errorf("failed to update data with id '%v' : %v", crumb.GetId(), err)
		return nil, err
	}

	return &pb.Id{
		Value: crumb.GetId(),
	}, nil
}

func (r *Route) Delete(ctx context.Context, id *pb.Id) (*pb.Id, error) {
	r.lc.Debug("received new Delete request")
	err := r.dbClient.Delete(r.dbConfig.Name, r.dbConfig.Collection, id.GetValue())
	if err != nil {
		r.lc.Errorf("failed to delete data with id '%v': %v", id.GetValue(), err)
		return nil, err
	}
	return id, nil
}
