package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Date struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Owner      primitive.ObjectID `json:"ownerId" bson:"ownerId"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Ethnicity  string             `json:"ethnicity,omitempty" bson:"ethnicity,omitempty"`
	Occupation string             `json:"occupation,omitempty" bson:"occupation,omitempty"`
	Place      string             `json:"place,omitempty" bson:"place,omitempty"`
	TypeOfDate string             `json:"typeOfDate" bson:"typeOfDate"`
	Cost       float32            `json:"cost,omitempty" bson:"cost,omitempty"`
	Result     string             `json:"result,omitempty" bson:"result,omitempty"`
	Age        int32              `json:"age" bson:"age"`
	Date       time.Time          `json:"date" bson:"date"`
	Split      bool               `json:"split,omitempty" bson:"split,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}
