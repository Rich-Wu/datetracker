package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Date struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Owner      primitive.ObjectID `json:"ownerId" bson:"ownerId"`
	FirstName  string             `json:"firstName" bson:"firstName" form:"first_name" binding:"required"`
	LastName   string             `json:"lastName,omitempty" bson:"lastName,omitempty" form:"last_name"`
	Ethnicity  string             `json:"ethnicity,omitempty" bson:"ethnicity,omitempty" form:"ethnicity" binding:"required"`
	Occupation string             `json:"occupation,omitempty" bson:"occupation,omitempty" form:"occupation"`
	Place      string             `json:"place,omitempty" bson:"place,omitempty" form:"place"`
	TypeOfDate string             `json:"typeOfDate" bson:"typeOfDate" form:"type_of_date" binding:"required"`
	Cost       float32            `json:"cost,omitempty" bson:"cost,omitempty,truncate" form:"cost"`
	Result     string             `json:"result,omitempty" bson:"result,omitempty" form:"result"`
	Age        int32              `json:"age" bson:"age" form:"age" binding:"required"`
	Date       time.Time          `json:"date" bson:"date" form:"date" binding:"required"`
	Split      bool               `json:"split,omitempty" bson:"split,omitempty" form:"split" binding:"required"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}
