package routes

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/go-playground/validator/v10"
	"github.com/haguru/horus/followerdb/config"
	pb "github.com/haguru/horus/followerdb/internal/routes/protos"
	grpcMocks "github.com/haguru/horus/followerdb/internal/routes/protos/mocks"
	"github.com/haguru/horus/followerdb/pkg/interfaces/mocks"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRoute_AddFollow(t *testing.T) {
	type args struct {
		ctx    context.Context
		follow *pb.Follow
	}
	tests := []struct {
		name         string
		args         args
		clientRtn    string
		clientErrRtn error
		want         *pb.Id
		wantErr      bool
	}{
		{
			name: "Successful  Addfollow",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "test_userid",
					FollowerId: "test_follower_userid",
				},
			},
			clientRtn:    "test_followid",
			clientErrRtn: nil,
			want:         &pb.Id{Value: "test_followid"},
			wantErr:      false,
		},
		{
			name: "validation fail",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "",
					FollowerId: "test_follower_userid",
				},
			},
			clientRtn:    "",
			clientErrRtn: nil,
			want:         nil,
			wantErr:      true,
		},
		{
			name: "client error",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "test_userid",
					FollowerId: "test_follower_userid",
				},
			},
			clientRtn:    "",
			clientErrRtn: fmt.Errorf("failed"),
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(tt.clientRtn, tt.clientErrRtn).Maybe()
			r := &Route{
				dbConfig: &config.Database{
					Name:       "test_database",
					Collection: "test_collection",
				},
				dbClient:  mockClient,
				lc:        logger.NewMockClient(),
				validator: validator.New(),
			}
			got, err := r.AddFollow(tt.args.ctx, tt.args.follow)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.AddFollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.AddFollow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_GetFollowers(t *testing.T) {
	type args struct {
		id *pb.Id
	}
	tests := []struct {
		name         string
		clientRtn    interface{}
		clientErrRtn error
		streamErrRtn error
		args         args
		wantErr      bool
	}{
		{
			name:         "Success GetFollowers",
			clientRtn:    []bson.D{{{Key: "test_userId", Value: "test_followerUserId"}, {Key: "test_userId", Value: "test_followerUserId2"}}},
			clientErrRtn: nil,
			streamErrRtn: nil,
			args: args{
				id: &pb.Id{
					Value: "test_userId",
				},
			},
			wantErr: false,
		},
		{
			name:         "client error",
			clientRtn:    nil,
			clientErrRtn: fmt.Errorf("failed"),
			streamErrRtn: nil,
			args: args{
				id: &pb.Id{
					Value: "test_userId",
				},
			},
			wantErr: true,
		},
		{
			name:         "stream error",
			clientRtn:    []bson.D{{{Key: "test_userId", Value: "test_followerUserId"}, {Key: "test_userId", Value: "test_followerUserId2"}}},
			clientErrRtn: nil,
			streamErrRtn: fmt.Errorf("failed"),
			args: args{
				id: &pb.Id{
					Value: "test_userId",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return(tt.clientRtn, tt.clientErrRtn).Maybe()
			streamServerMock := grpcMocks.NewServerStreamingServer[pb.Id](t)
			streamServerMock.On("Send", mock.Anything).Return(tt.streamErrRtn).Maybe()
			r := &Route{
				dbConfig: &config.Database{
					Name:       "test_database",
					Collection: "test_collection",
				},
				dbClient: mockClient,
				lc:       logger.NewMockClient(),
			}
			if err := r.GetFollowers(tt.args.id, streamServerMock); (err != nil) != tt.wantErr {
				t.Errorf("Route.GetFollowers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoute_Unfollow(t *testing.T) {
	type args struct {
		ctx    context.Context
		follow *pb.Follow
	}
	tests := []struct {
		name         string
		args         args
		clientErrRtn error
		want         *pb.Status
		wantErr      bool
	}{
		{
			name: "successful  unfollow",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "test_userId",
					FollowerId: "test_followerUserId",
				},
			},
			clientErrRtn: nil,
			want: &pb.Status{
				Value: 200,
			},
			wantErr: false,
		},
		{
			name: "validation error",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "",
					FollowerId: "test_followerUserId",
				},
			},
			clientErrRtn: nil,
			want:         nil,
			wantErr:      true,
		},
		{
			name: "client error",
			args: args{
				ctx: context.Background(),
				follow: &pb.Follow{
					Id:         "test_userId",
					FollowerId: "test_followerUserId",
				},
			},
			clientErrRtn: fmt.Errorf("failed"),
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewDbClient(t)
			mockClient.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(tt.clientErrRtn).Maybe()
			r := &Route{
				dbConfig: &config.Database{
					Name:       "test_database",
					Collection: "test_collection",
				},
				dbClient:  mockClient,
				lc:        logger.NewMockClient(),
				validator: validator.New(),
			}
			got, err := r.Unfollow(tt.args.ctx, tt.args.follow)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Unfollow() = %v, want %v", got, tt.want)
			}
		})
	}
}
