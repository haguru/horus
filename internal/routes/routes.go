package routes

import (
	"context"

	"github.com/haguru/horus/config"
	pb "github.com/haguru/horus/internal/routes/protos"
	"github.com/haguru/horus/pkg/models"
	"github.com/haguru/horus/pkg/mongodb/interfaces"
	dbModel "github.com/haguru/horus/pkg/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Route struct {
	dbCconfig *config.Database
	dbClient  interfaces.Client
}

func (r *Route) Create(ctx context.Context, crumb *pb.Crumb) (*pb.Id, error) {
	c := models.Crumb{
		Location: dbModel.Point{
			Type:        crumb.GetLocation().GetType(),
			Coordinates: crumb.GetLocation().GetCoordinates(),
		},
		Message: crumb.GetMessagecrumb(),
		User:    crumb.User,
	}
	id, err := r.dbClient.InsertRecord(r.dbCconfig.Name, r.dbCconfig.Collection, c)
	if err != nil {
		return nil, err
	}

	return &pb.Id{Value: id}, nil
}

func (r *Route) GetCrumbs(point *pb.Point, stream pb.CrumbDB_GetCrumbsServer) error {
	p := dbModel.Point{
		Type:        point.GetType(),
		Coordinates: point.GetCoordinates(),
	}
	data, err := r.dbClient.SpaitalQuery(p, r.dbCconfig.Name, r.dbCconfig.Collection)
	if err != nil {
		return err
	}
	for _, item := range data {
		// unmarshall data to grpc data type
		doc, err := bson.Marshal(item)
		if err != nil {
			return err
		}

		crumb := &pb.Crumb{}
		err = bson.Unmarshal(doc, crumb)
		if err != nil {
			return err
		}

		// send crumb
		err = stream.Send(crumb)
		if err != nil {
			return err
		}

	}
	return nil
}

func (r *Route) Update(crumb *pb.Crumb) (*pb.Id, error) {
	c := models.MessageUpdateRequest{
		Message: crumb.GetMessagecrumb(),
	}
	err := r.dbClient.Update(r.dbCconfig.Name, r.dbCconfig.Collection, crumb.Id, c)
	if err != nil {
		return nil, err
	}

	return &pb.Id{
		Value: crumb.GetId(),
	}, nil
}

func (r *Route) Delete(id *pb.Id) (*pb.Id, error) {
	err := r.dbClient.Delete(r.dbCconfig.Name, r.dbCconfig.Collection, id.GetValue())
	if err != nil {
		return nil, err
	}
	return id, nil
}
