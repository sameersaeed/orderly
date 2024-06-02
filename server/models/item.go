package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Item struct {
    ID          bson.ObjectId `bson:"_id,omitempty"`
    Name        string        `bson:"name"`
    Description string        `bson:"description"`
    Price       float64       `bson:"price"`
}