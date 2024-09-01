package models

import (
	pb "github.com/haguru/horus/crumbdb/internal/routes/protos"
)

type MessageUpdateRequest struct {
	Message string `bson:"message" json:"message"`
}

type GrpcServer struct {
	pb.UnimplementedCrumbDBServer
}
