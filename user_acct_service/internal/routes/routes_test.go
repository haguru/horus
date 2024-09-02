package routes

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/haguru/horus/useracctdb/config"
	pb "github.com/haguru/horus/useracctdb/internal/routes/protos"
	"github.com/haguru/horus/useracctdb/pkg/interfaces"
	"github.com/haguru/horus/useracctdb/pkg/interfaces/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRoute_Create(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}
	type args struct {
		ctx  context.Context
		user *pb.User
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		dbClientRtn error
		want        *pb.Id
		wantErr     bool
	}{
		{
			name: "sucessful create",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Id:       "test_id",
					Email:    "test@horus.com",
					Username: "test_username",
					Password: "test_password",
				},
			},
			dbClientRtn: nil,
			want: &pb.Id{
				Value: "test_id",
			},
			wantErr: false,
		},
		{
			name: "validation error - no username",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Id:       "test_id",
					Email:    "test@horus.com",
					Username: "",
					Password: "test_password",
				},
			},
			dbClientRtn: fmt.Errorf("no username given"),
			want:        nil,
			wantErr:     true,
		},
		{
			name: "validation error - no password",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Id:       "test_id",
					Email:    "test@horus.com",
					Username: "test_username",
					Password: "",
				},
			},
			dbClientRtn: fmt.Errorf("no password given"),
			want:        nil,
			wantErr:     true,
		},
		{
			name: "client error",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Id:       "test_id",
					Email:    "test@horus.com",
					Username: "test_username",
					Password: "test_password",
				},
			},
			dbClientRtn: fmt.Errorf("client failed"),
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(tt.dbClientRtn)
			r := &Route{
				dbCconfig: tt.fields.dbCconfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
			}
			got, err := r.Create(tt.args.ctx, tt.args.user)
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

func TestRoute_GetUser(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}
	type args struct {
		ctx     context.Context
		userReq *pb.UserRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		dbCLientRtn error
		want        *pb.User
		wantErr     bool
	}{
		{
			name: "successful get user",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				userReq: &pb.UserRequest{
					Username: "test_username",
				},
			},
			dbCLientRtn: nil,
			want: &pb.User{
				Id:       "test_id",
				Email:    "test@horus.com",
				Username: "test_username",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name: "validation error- no username",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				userReq: &pb.UserRequest{
					Username: "",
				},
			},
			dbCLientRtn: fmt.Errorf("no username given"),
			want:        nil,
			wantErr:     true,
		},
		{
			name: "client error",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				userReq: &pb.UserRequest{
					Username: "test_username",
				},
			},
			dbCLientRtn: fmt.Errorf("client failed"),
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("GetUser", mock.Anything).Return(tt.dbCLientRtn)
			r := &Route{
				dbCconfig: tt.fields.dbCconfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
			}
			got, err := r.GetUser(tt.args.ctx, tt.args.userReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestRoute_Update(t *testing.T) {
	type fields struct {
		dbCconfig                     *config.Database
		dbClient                      interfaces.DbClient
		lc                            logger.LoggingClient
		UnimplementedUserAcctDBServer pb.UnimplementedUserAcctDBServer
	}
	type args struct {
		ctx       context.Context
		passwdReq *pb.PasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Status
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				dbCconfig:                     tt.fields.dbCconfig,
				dbClient:                      tt.fields.dbClient,
				lc:                            tt.fields.lc,
				UnimplementedUserAcctDBServer: tt.fields.UnimplementedUserAcctDBServer,
			}
			got, err := r.Update(tt.args.ctx, tt.args.passwdReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestRoute_Delete(t *testing.T) {
	type fields struct {
		dbCconfig                     *config.Database
		dbClient                      interfaces.DbClient
		lc                            logger.LoggingClient
		UnimplementedUserAcctDBServer pb.UnimplementedUserAcctDBServer
	}
	type args struct {
		ctx     context.Context
		userReq *pb.UserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Status
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				dbCconfig:                     tt.fields.dbCconfig,
				dbClient:                      tt.fields.dbClient,
				lc:                            tt.fields.lc,
				UnimplementedUserAcctDBServer: tt.fields.UnimplementedUserAcctDBServer,
			}
			got, err := r.Delete(tt.args.ctx, tt.args.userReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
