package routes

import (
	"context"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/haguru/horus/crumbdb/config"
	pb "github.com/haguru/horus/crumbdb/internal/routes/protos"
	"github.com/haguru/horus/crumbdb/pkg/models"
	"github.com/haguru/horus/crumbdb/pkg/mongodb/interfaces"
	"go.mongodb.org/mongo-driver/bson"
)

type Route struct {
	dbCconfig *config.Database
	dbClient  interfaces.Client
	lc        logger.LoggingClient
	pb.UnimplementedCrumbDBServer
}

func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.Client) *Route {
	return &Route{
		dbCconfig: config,
		dbClient:  dbclient,
		lc:        lc,
	}
}

func (r *Route) Create(ctx context.Context, crumb *pb.Crumb) (*pb.Id, error) {
	r.lc.Debugf("received Create request: %v", crumb)

	id, err := r.dbClient.InsertRecord(r.dbCconfig.Name, r.dbCconfig.Collection, crumb)
	if err != nil {
		return nil, err
	}

	return &pb.Id{Value: id}, nil
}

func (r *Route) GetCrumbs(point *pb.Point, stream pb.CrumbDB_GetCrumbsServer) error {
	r.lc.Debug("received new GetCumbs request")

	data, err := r.dbClient.SpaitalQuery(point, r.dbCconfig.Name, r.dbCconfig.Collection)
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

	messageUpdate := models.MessageUpdateRequest{
		Message: crumb.GetMessage(),
	}

	err := r.dbClient.Update(r.dbCconfig.Name, r.dbCconfig.Collection, crumb.Id, messageUpdate)
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
	err := r.dbClient.Delete(r.dbCconfig.Name, r.dbCconfig.Collection, id.GetValue())
	if err != nil {
		r.lc.Errorf("failed to delete data with id '%v': %v", id.GetValue(), err)
		return nil, err
	}
	return id, nil
}
