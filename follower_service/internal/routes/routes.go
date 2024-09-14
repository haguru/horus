package routes

import (
	"context"
	"fmt"

	"github.com/haguru/horus/followerdb/config"
<<<<<<< HEAD
	"github.com/haguru/horus/followerdb/pkg/interfaces"
	appMetrics "github.com/haguru/horus/followerdb/pkg/prometheus"
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	
	
=======
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"
	"github.com/haguru/horus/followerdb/pkg/interfaces"
	appMetrics "github.com/haguru/horus/followerdb/pkg/prometheus"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
>>>>>>> 5d11d8811ccd7f1f01c900e91faea8139cd8b878
)

type Route struct {
	dbConfig  *config.Database
	dbClient  interfaces.DbClient
	lc        logger.LoggingClient
	validator *validator.Validate
	metrics   *appMetrics.Metrics
	pb.UnimplementedFollowerDBServer
}

// TODO
func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.DbClient, validator *validator.Validate, metrics *appMetrics.Metrics) *Route {
	return &Route{
		dbConfig:  config,
		dbClient:  dbclient,
		lc:        lc,
		validator: validator,
<<<<<<< HEAD
		metrics: metrics,
=======
		metrics:   metrics,
>>>>>>> 5d11d8811ccd7f1f01c900e91faea8139cd8b878
	}
}

func (r *Route) AddFollow(ctx context.Context, follow *pb.Follow) (*pb.Id, error) {
	r.lc.Debugf("received AddFollow request")
	// r.metrics.RequestsCount.Inc()

	// Validate the User struct
	err := r.validator.Struct(follow)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	id, err := r.dbClient.Create(r.dbConfig.Name, r.dbConfig.Collection, follow)
	if err != nil {
		return nil, fmt.Errorf("failed to add follow: %v", err)
	}

	return &pb.Id{Value: id}, nil
}

func (r *Route) GetFollowers(id *pb.Id, stream pb.FollowerDB_GetFollowersServer) error {
	r.lc.Debugf("received Getfollowers request")
	// r.metrics.RequestsCount.Inc()

	filter := map[string]interface{}{"userId": id.GetValue()}
	items, err := r.dbClient.GetAll(r.dbConfig.Name, r.dbConfig.Collection, filter)
	if err != nil {
		return fmt.Errorf("failed to retrieve follows for id %v: %v", id.GetValue(), err)
	}
	for _, item := range items.([]bson.D) {
		// unmarshall data to grpc data type
		doc, err := bson.Marshal(item)
		if err != nil {
			r.lc.Errorf("failed to marshal an item in data: %v", err)
			return err
		}

		follow := &pb.Follow{}
		err = bson.Unmarshal(doc, follow)
		if err != nil {
			r.lc.Errorf("failed to unmarshal an item in data: %v", err)
			return err
		}

		// send crumb
		id := &pb.Id{Value: follow.GetFollowerId()}
		err = stream.Send(id)
		if err != nil {
			r.lc.Errorf("failed to send item in data: %v", err)
			return err
		}
	}
	return nil
}

func (r *Route) Unfollow(ctx context.Context, follow *pb.Follow) (*pb.Status, error) {
	r.lc.Debugf("received Unfollow request")
	// r.metrics.RequestsCount.Inc()

	// Validate the User struct
	err := r.validator.Struct(follow)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	filter := map[string]interface{}{"userId": follow.GetId(), "followerUserId": follow.GetFollowerId()}
	err = r.dbClient.Delete(r.dbConfig.Name, r.dbConfig.Collection, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete follow: %v", err)
	}

	// do I reall want to return a status?
	return &pb.Status{Value: 200}, nil
}
