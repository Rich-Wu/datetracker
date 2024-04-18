package main

import "time"

type Email struct {
	Address    string    `bson:"email" json:"email" binding:"required,email" form:"email"`
	SignupTime time.Time `bson:"signupTime" json:"signupTime"`
}
