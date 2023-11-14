package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Medication struct {
	Id   primitive.ObjectID `json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
}
