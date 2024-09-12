package routes

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/go-playground/validator/v10"
	"github.com/haguru/horus/crumbdb/config"
	pb "github.com/haguru/horus/crumbdb/internal/routes/protos"
	grpcMock "github.com/haguru/horus/crumbdb/internal/routes/protos/mocks"
	"github.com/haguru/horus/crumbdb/pkg/mongodb"
	"github.com/haguru/horus/crumbdb/pkg/mongodb/interfaces/mocks"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRoute_Create(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}
	type args struct {
		ctx   context.Context
		crumb *pb.Crumb
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		insertRecordRtn string
		errorRtn        error
		want            *pb.Id
		wantErr         bool
	}{
		{
			name: "success full create",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.MockLogger{},
			},
			args: args{
				ctx: context.Background(),
				crumb: &pb.Crumb{
					Location: &pb.Point{
						Coordinates: []float64{-122.66025176499872, 45.692956992343845},
						Type:        "Point",
					},
					User:    "test_user",
					Message: "this is simply a test",
				},
			},
			insertRecordRtn: "test_id",
			errorRtn:        nil,
			want: &pb.Id{
				Value: "test_id",
			},
			wantErr: false,
		},
		{
			name: "fail to create",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.MockLogger{},
			},
			args: args{
				ctx:   context.Background(),
				crumb: &pb.Crumb{
					Location: &pb.Point{
						Type: mongodb.POINT_TYPE_POINT,
						Coordinates: []float64{-122.66025176499872, 45.692956992343845},
					},
					User: "test_user",
					Message: "test_message",
				},
			},
			insertRecordRtn: "",
			errorRtn:        fmt.Errorf("failed"),
			want:            nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewClient(t)
			mockClient.On("InsertRecord", mock.Anything, mock.Anything, mock.Anything).Return(tt.insertRecordRtn, tt.errorRtn)

			r := &Route{
				dbConfig: tt.fields.dbCconfig,
				dbClient: mockClient,
				lc:       tt.fields.lc,
				validator: validator.New(),
			}

			got, err := r.Create(tt.args.ctx, tt.args.crumb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Value, tt.want.Value) {
				t.Errorf("Route.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_GetCrumbs(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}

	tests := []struct {
		name           string
		fields         fields
		streamErrRtn   error
		clientRtn      []bson.D
		clientErrorRtn error
		point          *pb.Point
		wantErr        bool
	}{
		{
			name: "succesfully get list of crumbs",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			streamErrRtn:   nil,
			clientRtn:      []bson.D{{{Key: "user", Value: "test"}}},
			clientErrorRtn: nil,
			point: &pb.Point{
				Coordinates: []float64{-122.66025176499872, 45.692956992343845},
				Type:        "Point",
			},
			wantErr: false,
		},
		{
			name: "stream fail to send list of crumbs",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			streamErrRtn:   fmt.Errorf("failed"),
			clientRtn:      []bson.D{{{Key: "user", Value: "test_user"}}},
			clientErrorRtn: nil,
			point: &pb.Point{
				Coordinates: []float64{-122.66025176499872, 45.692956992343845},
				Type:        "Point",
			},
			wantErr: true,
		},
		{
			name: "client fail to get list of crumbs",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			streamErrRtn:   nil,
			clientRtn:      []bson.D{},
			clientErrorRtn: fmt.Errorf("failed"),
			point: &pb.Point{
				Coordinates: []float64{-122.66025176499872, 45.692956992343845},
				Type:        "Point",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := grpcMock.NewServerStreamingServer[pb.Crumb](t)
			stream.On("Send", mock.Anything).Return(tt.streamErrRtn).Maybe()
			mockClient := mocks.NewClient(t)
			mockClient.On("SpaitalQuery", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.clientRtn, tt.clientErrorRtn)
			r := &Route{
				dbConfig: tt.fields.dbCconfig,
				dbClient: mockClient,
				lc:       tt.fields.lc,
				validator: validator.New(),
			}
			if err := r.GetCrumbs(tt.point, stream); (err != nil) != tt.wantErr {
				t.Errorf("Route.GetCrumbs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoute_Update(t *testing.T) {
	type fields struct {
		dbCconfig *config.Database
		lc        logger.LoggingClient
	}
	type args struct {
		ctx   context.Context
		crumb *pb.Crumb
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		clientRtn      string
		clientErrorRtn error
		want           *pb.Id
		wantErr        bool
	}{
		{
			name: "successful update",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				crumb: &pb.Crumb{
					Id:      "test_id",
					Message: "testing update method",
				},
			},
			clientErrorRtn: nil,
			want: &pb.Id{
				Value: "test_id",
			},
			wantErr: false,
		},
		{
			name: "client failed to update",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				crumb: &pb.Crumb{
					Id:      "test_id",
					Message: "testing update method",
				},
			},
			clientErrorRtn: fmt.Errorf("fail"),
			want:           nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewClient(t)
			mockClient.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.clientErrorRtn)
			r := &Route{
				dbConfig: tt.fields.dbCconfig,
				dbClient: mockClient,
				lc:       tt.fields.lc,
				validator: validator.New(),
			}
			got, err := r.Update(tt.args.ctx, tt.args.crumb)
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
		ctx context.Context
		id  *pb.Id
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		clientErrorRtn error
		want           *pb.Id
		wantErr        bool
	}{
		{
			name: "succesful delete",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				id: &pb.Id{
					Value: "test_id",
				},
			},
			clientErrorRtn: nil,
			want: &pb.Id{
				Value: "test_id",
			},
			wantErr: false,
		},
		{
			name: "succesful delete",
			fields: fields{
				dbCconfig: &config.Database{
					Name:       "test",
					Collection: "test",
				},
				lc: logger.NewMockClient(),
			},
			args: args{
				ctx: context.Background(),
				id: &pb.Id{
					Value: "test_id",
				},
			},
			clientErrorRtn: fmt.Errorf("failed"),
			want:           nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewClient(t)
			mockClient.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(tt.clientErrorRtn)
			r := &Route{
				dbConfig: tt.fields.dbCconfig,
				dbClient: mockClient,
				lc:       tt.fields.lc,
				validator: validator.New(),
			}
			got, err := r.Delete(tt.args.ctx, tt.args.id)
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
