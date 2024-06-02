package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order_submission_tool/bot"
	"order_submission_tool/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// submitting order to discord channel for processing
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    orderDetails := fmt.Sprintf("**New order received! @here**\n\nUser ID: %s\nUser email: %s\n```Items ordered:\n", order.UserID, order.Email)
    for _, item := range order.Cart {
        orderDetails +=  " *  " + item.ItemName + " (" + fmt.Sprintf("%.2f x%d", item.Price, item.Quantity) + ")\n"
    }
    orderDetails += fmt.Sprintf("```Total price: %.2f", order.TotalPrice)

    // sending message with order details through discord bot
    err = bot.SendMessage(orderDetails)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // saving order to db
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("orders")
    _, err = collection.InsertOne(ctx, order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Print("LOG: successfully created order: ", order.OrderID)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "the order has been sent to Discord channel successfully"})
}

// fetching past orders to display to user
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := db.Collection("orders")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var orders []models.Order
    for cursor.Next(ctx) {
        var order models.Order
        
        err := cursor.Decode(&order)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        orders = append(orders, order)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(orders)
}