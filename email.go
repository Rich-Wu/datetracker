package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name" json:"name" form:"name" binding:"required"`
	Address    string             `bson:"email" json:"email" form:"email" binding:"required,email"`
	SignupTime time.Time          `bson:"signupTime" json:"signupTime"`
}
