package routes

import (
	"context"
	"fmt"

	"github.com/haguru/horus/follower_service/config"
	pb "github.com/haguru/horus/follower_service/internal/routes/protos"
	"github.com/haguru/horus/follower_service/pkg/interfaces"
	appMetrics "github.com/haguru/horus/follower_service/pkg/prometheus"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	reqValidator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type Route struct {
	dbConfig  *config.Database
	dbClient  interfaces.DbClient
	lc        logger.LoggingClient
	validator *reqValidator.Validate
	metrics   *appMetrics.Metrics
	pb.UnimplementedFollowerDBServer
}

// NewRoute creates a new Route instance with the provided dependencies
func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.DbClient, validator *reqValidator.Validate) *Route {
	// Create a new Route instance with the provided dependencies
	// Note: metrics will be injected later by the app.go when initializing the server
	return &Route{
		dbConfig:  config,
		dbClient:  dbclient,
		lc:        lc,
		validator: validator,
	}
}

// SetMetrics sets the metrics for the Route
func (r *Route) SetMetrics(metrics *appMetrics.Metrics) {
	r.metrics = metrics
}

func (r *Route) AddFollow(ctx context.Context, follow *pb.Follow) (*pb.Id, error) {
	r.lc.Debugf("received AddFollow request")
	// Metrics are automatically collected by the gRPC interceptors
	if r.metrics != nil {
		r.metrics.FollowsAdded.Inc()
	}

	// Validate the User struct
	err := r.validator.Struct(follow)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(reqValidator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	id, err := r.dbClient.Create(r.dbConfig.DatabaseName, r.dbConfig.Collection, follow)
	if err != nil {
		return nil, fmt.Errorf("failed to add follow: %v", err)
	}

	return &pb.Id{Value: id}, nil
}

func (r *Route) GetFollowers(req *pb.GetFollowersRequest, stream pb.FollowerDB_GetFollowersServer) error {
	r.lc.Debugf("received Getfollowers request")
	// Metrics are automatically collected by the gRPC interceptors
	if r.metrics != nil {
		r.metrics.FollowersRetrieved.Inc()
	}

	filter := map[string]interface{}{"userId": req.GetUserId()}
	skip := (req.GetPage() - 1) * req.GetPageSize()
	limit := req.GetPageSize()
	items, err := r.dbClient.GetAll(r.dbConfig.DatabaseName, r.dbConfig.Collection, filter, skip, limit)
	if err != nil {
		return fmt.Errorf("failed to retrieve follows for id %v: %v", req.GetUserId(), err)
	}
	r.lc.Debugf("items: %v", items)
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
		r.lc.Debugf("stream.Send success")
	}
	return nil
}

func (r *Route) Unfollow(ctx context.Context, follow *pb.Follow) (*pb.Status, error) {
	r.lc.Debugf("received Unfollow request")
	// Metrics are automatically collected by the gRPC interceptors
	if r.metrics != nil {
		r.metrics.UnfollowsProcessed.Inc()
	}

	// Validate the User struct
	err := r.validator.Struct(follow)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(reqValidator.ValidationErrors)

		return nil, fmt.Errorf("validation error: %s", errors)
	}

	filter := map[string]interface{}{"userId": follow.GetId(), "followerUserId": follow.GetFollowerId()}
	err = r.dbClient.Delete(r.dbConfig.DatabaseName, r.dbConfig.Collection, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete follow: %v", err)
	}

	// Return a status with HTTP 200 OK to indicate success
	return &pb.Status{Value: 200}, nil
}
