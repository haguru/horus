package routes

import (
	"context"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/haguru/horus/config"
	pb "github.com/haguru/horus/internal/routes/protos"
	"github.com/haguru/horus/pkg/mongodb/interfaces"
)

func TestRoute_Create(t *testing.T) {
	type fields struct {
		dbCconfig                  *config.Database
		dbClient                   interfaces.Client
		lc                         logger.LoggingClient
		UnimplementedCrumbDBServer pb.UnimplementedCrumbDBServer
	}
	type args struct {
		ctx   context.Context
		crumb *pb.Crumb
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Id
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				dbCconfig:                  tt.fields.dbCconfig,
				dbClient:                   tt.fields.dbClient,
				lc:                         tt.fields.lc,
				UnimplementedCrumbDBServer: tt.fields.UnimplementedCrumbDBServer,
			}
			got, err := r.Create(tt.args.ctx, tt.args.crumb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
