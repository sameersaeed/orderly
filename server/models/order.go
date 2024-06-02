package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
    ID       	bson.ObjectId 	`bson:"_id,omitempty"`
    Email    	string        	`json:"email"`
    Cart        []CartItem      `json:"cart"`
    TotalPrice 	float64       	`json:"total_price"`
    Message     string        	`json:"message,omitempty"`
}

type CartItem struct {
    ItemID   bson.ObjectId `bson:"item_id,omitempty"`
    Quantity int           `json:"quantity"`
}