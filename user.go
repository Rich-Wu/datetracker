package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserName  string             `json:"username" bson:"username" form:"username" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName" form:"first_name" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName" form:"last_name"`
	Password  string             `bson:"password,-" form:"password" binding:"required"`
}
