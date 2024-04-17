package main

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserName   string             `json:"username" bson:"username" form:"username" binding:"required"`
	FirstName  string             `json:"firstName" bson:"firstName" form:"first_name" binding:"required"`
	LastName   string             `json:"lastName" bson:"lastName" form:"last_name"`
	Password   string             `bson:"password,-" json:"password,omitempty" form:"password" binding:"required"`
	ProfilePic string             `json:"profilePic,omitempty" bson:"profilePic,omitempty" form:"profilePic"`
}

type ProfilePic struct {
	Owner    primitive.ObjectID `bson:"owner" json:"owner"`
	FileName string             `bson:"fileName" json:"fileName"`
	Data     primitive.Binary   `bson:"data" json:"data"`
	Metadata ProfilePicMetadata `bson:"metadata" json:"metadata"`
}

type ProfilePicMetadata struct {
	Filetype string            `bson:"filetype" json:"filetype"`
	Tags     map[string]string `bson:"tags" json:"tags"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Password = ""
	return json.Marshal(x)
}
