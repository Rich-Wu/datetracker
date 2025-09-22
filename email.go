package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name" json:"name" form:"name" binding:"required"`
	Address     string             `bson:"email" json:"email" form:"email" binding:"required,email"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber" form:"phonenumber" binding:"omitempty,len=10,numeric"`
	Sms         bool               `bson:"sms" json:"sms" form:"sms"`
	SignupTime  time.Time          `bson:"signupTime" json:"signupTime"`
}
