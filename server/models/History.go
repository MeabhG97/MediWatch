package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type History struct {
	Id          primitive.ObjectID   `json:"id,omitempty"`
	DateTime    primitive.DateTime   `json:"datetime,omitempty"`
	Taken       bool                 `json:"taken"`
	Medications []primitive.ObjectID `json:"medications,omitempty"`
}
