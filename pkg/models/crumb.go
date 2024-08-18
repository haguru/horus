package models

import (
	"github.com/haguru/horus/pkg/mogodb/models"
)

type Crumb struct {
	Location models.Point
	Message  string `bson:"message,omitempty" json:"message,omitempty"`
	User     string `bson:"user,omitempty" json:"user,omitempty"`
}
