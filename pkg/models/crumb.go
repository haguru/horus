package models

import (
	"github.com/haguru/horus/pkg/mongodb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Crumb struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Location models.Point       `bson:"location" json:"location"`
	Message  string             `bson:"message,omitempty" json:"messagecrumb,omitempty"`
	User     string             `bson:"user,omitempty" json:"user,omitempty"`
}

type MessageUpdateRequest struct {
	Message string `bson:"message,omitempty" json:"messagecrumb,omitempty"`
}
