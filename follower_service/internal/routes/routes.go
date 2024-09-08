package routes

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/followerdb/config"
	"github.com/haguru/horus/followerdb/pkg/interfaces"
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"
)

// TODO
type Route struct{

	pb.UnimplementedFollowerDBServer
}

// TODO
func NewRoute(lc logger.LoggingClient, config *config.Database, dbclient interfaces.DbClient) *Route {
	return nil
}
