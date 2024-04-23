package main

import "time"

type Email struct {
	Name       string    `bson:"name" json:"name" form:"name" binding:"required"`
	Address    string    `bson:"email" json:"email" form:"email" binding:"required,email"`
	SignupTime time.Time `bson:"signupTime" json:"signupTime"`
}
