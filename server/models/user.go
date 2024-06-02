package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
    ID       bson.ObjectId `bson:"_id,omitempty"`
    Name     string        `bson:"name"`
    Email    string        `bson:"email"`
    Password string        `bson:"password"`
    IsAdmin  bool          `bson:"isAdmin"`
    Token    string        `bson:"token"` 
}