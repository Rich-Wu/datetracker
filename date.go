package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Date struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId    primitive.ObjectID `json:"ownerId" bson:"ownerId"`
	FirstName  string             `json:"firstName" bson:"firstName" form:"first_name" binding:"required"`
	LastName   string             `json:"lastName,omitempty" bson:"lastName,omitempty" form:"last_name"`
	Age        int32              `json:"age" bson:"age" form:"age" binding:"required"`
	Occupation string             `json:"occupation,omitempty" bson:"occupation,omitempty" form:"occupation"`
	Ethnicity  string             `json:"ethnicity,omitempty" bson:"ethnicity,omitempty" form:"ethnicity" binding:"required"`
	Places     []*Place           `json:"places" bson:"places" form:"places" binding:"required,min=1"`
	Cost       float32            `json:"cost" bson:"cost,truncate" binding:"required"`
	Result     string             `json:"result,omitempty" bson:"result,omitempty" form:"result"`
	Date       time.Time          `json:"date" bson:"date" form:"date" binding:"required"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

type Place struct {
	Place       string  `json:"place,omitempty" bson:"place,omitempty" form:"place" bind:"required"`
	TypeOfPlace string  `json:"typeOfPlace" bson:"typeOfPlace" form:"type_of_place" binding:"required"`
	Cost        float32 `json:"cost,omitempty" bson:"cost,omitempty,truncate" form:"cost"`
	Split       bool    `json:"split,omitempty" bson:"split,omitempty" form:"split" binding:"required,oneof=True False"`
}
