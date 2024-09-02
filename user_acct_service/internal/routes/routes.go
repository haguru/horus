package routes

import (
	"context"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/useracctdb/config"
	"github.com/haguru/horus/useracctdb/pkg/interfaces"

	pb "github.com/haguru/horus/useracctdb/internal/routes/protos"
)

type Route struct {
	dbCconfig *config.Database
	dbClient  interfaces.DbClient
	lc        logger.LoggingClient
	pb.UnimplementedUserAcctDBServer
}

func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.DbClient) *Route {
	return &Route{
		dbCconfig: config,
		dbClient:  dbclient,
		lc:        lc,
	}
}

// TODO
func (r *Route) Create(ctx context.Context, user *pb.User) (*pb.Id, error) {
	id := &pb.Id{}

	return id, nil
}

// TODO
func (r *Route) GetUser(ctx context.Context, userReq *pb.UserRequest) (*pb.User, error) {
	user := &pb.User{}

	return user, nil
}

// TODO
func (r *Route) Update(ctx context.Context, passwdReq *pb.PasswordRequest) (*pb.Status, error) {
	status := &pb.Status{}

	return status, nil
}

// TODO
func (r *Route) Delete(ctx context.Context, userReq *pb.UserRequest) (*pb.Status, error) {
	status := &pb.Status{}

	return status, nil
}
