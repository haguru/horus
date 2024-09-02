package routes

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	"github.com/haguru/horus/useracctdb/config"
	pb "github.com/haguru/horus/useracctdb/internal/routes/protos"
	"github.com/haguru/horus/useracctdb/pkg/interfaces/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRoute_Create(t *testing.T) {
	type fields struct {
		dbConfig *config.Database
		lc       logger.LoggingClient
	}
	type args struct {
		ctx  context.Context
		user *pb.User
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		dbClientRtn   error
		dbClientIDRtn string
		want          *pb.Id
		wantErr       bool
	}{
		{
			name: "sucessful create",
			fields: fields{
				dbConfig: &config.Database{
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
			dbClientRtn:   nil,
			dbClientIDRtn: "test_id",
			want: &pb.Id{
				Value: "test_id",
			},
			wantErr: false,
		},
		{
			name: "validation error - no username",
			fields: fields{
				dbConfig: &config.Database{
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
			dbClientRtn:   nil,
			dbClientIDRtn: "",
			want:          nil,
			wantErr:       true,
		},
		{
			name: "validation error - no password",
			fields: fields{
				dbConfig: &config.Database{
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
			dbClientRtn:   nil,
			dbClientIDRtn: "",
			want:          nil,
			wantErr:       true,
		},
		{
			name: "client error",
			fields: fields{
				dbConfig: &config.Database{
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
			dbClientRtn:   fmt.Errorf("client failed"),
			dbClientIDRtn: "",
			want:          nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Create", mock.Anything, mock.Anything, tt.args.user).Return(tt.dbClientIDRtn, tt.dbClientRtn).Maybe()
			r := &Route{
				dbConfig:  tt.fields.dbConfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
				validator: validator.New(),
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
		dbUserRtn   interface{}
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
					Email: "test@horus.com",
				},
			},
			dbCLientRtn: nil,
			dbUserRtn: pb.User{
				Id:       "test_id",
				Email:    "test@horus.com",
				Username: "test_username",
				Password: "test_password",
			},
			want: &pb.User{
				Id:       "test_id",
				Email:    "test@horus.com",
				Username: "test_username",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name: "validation error- no email",
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
					Email: "",
				},
			},
			dbCLientRtn: fmt.Errorf("failed"),
			dbUserRtn:   nil,
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
					Email: "test@horus.com",
				},
			},
			dbCLientRtn: fmt.Errorf("client failed"),
			dbUserRtn:   nil,
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(tt.dbUserRtn, tt.dbCLientRtn).Maybe()
			r := &Route{
				dbConfig:  tt.fields.dbCconfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
				validator: validator.New(),
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

func TestRoute_UpdatePassword(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}
	type args struct {
		ctx       context.Context
		passwdReq *pb.PasswordRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		dbClientRtn error
		want        *pb.Status
		wantErr     bool
	}{
		{
			name: "successful update",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				passwdReq: &pb.PasswordRequest{
					Email:    "test@horus.com",
					Password: "test_password",
				},
			},
			dbClientRtn: nil,
			want: &pb.Status{
				Value: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "dbclient fails",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				passwdReq: &pb.PasswordRequest{
					Email:    "test@horus.com",
					Password: "test_password",
				},
			},
			dbClientRtn: fmt.Errorf("failed"),
			want: &pb.Status{
				Value: http.StatusInternalServerError,
			},
			wantErr: true,
		},
		{
			name: "validation error - no email",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "horus",
					Collection: "users",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				passwdReq: &pb.PasswordRequest{
					Email:    "",
					Password: "test_password",
				},
			},
			dbClientRtn: nil,
			want: &pb.Status{
				Value: http.StatusBadRequest,
			},
			wantErr: true,
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
				passwdReq: &pb.PasswordRequest{
					Email:    "test@horus.com",
					Password: "",
				},
			},
			dbClientRtn: nil,
			want: &pb.Status{
				Value: http.StatusBadRequest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.dbClientRtn).Maybe()
			r := &Route{
				dbConfig:  tt.fields.dbCconfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
				validator: validator.New(),
			}
			got, err := r.UpdatePassword(tt.args.ctx, tt.args.passwdReq)
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

func TestRoute_Delete(t *testing.T) {
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
		want        *pb.Status
		wantErr     bool
	}{
		{
			name: "successfull delete",
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
					Email: "test@horus.com",
				},
			},
			dbCLientRtn: nil,
			want: &pb.Status{
				Value: http.StatusOK,
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
					Email: "",
				},
			},
			dbCLientRtn: nil,
			want: &pb.Status{
				Value: http.StatusBadRequest,
			},
			wantErr: true,
		},
		{
			name: "db client error",
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
					Email: "test@horus.com",
				},
			},
			dbCLientRtn: fmt.Errorf("failed"),
			want: &pb.Status{
				Value: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(tt.dbCLientRtn).Maybe()
			r := &Route{
				dbConfig:  tt.fields.dbCconfig,
				dbClient:  mockClient,
				lc:        tt.fields.lc,
				validator: validator.New(),
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
