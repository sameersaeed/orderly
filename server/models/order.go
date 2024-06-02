package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
    OrderID     primitive.ObjectID 	`bson:"_id,omitempty" json:"orderID,omitempty"`
    Email    	string        	    `json:"email"`
    Cart        []CartItems         `json:"cart"`
    UserID      primitive.ObjectID  `bson:"userID,omitempty" json:"userID,omitempty"`
    TotalPrice  float64             `json:"totalPrice"`
}

type CartItems struct {
    ItemName    string  `json:"item"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
}
