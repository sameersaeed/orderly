package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order_submission_tool/bot"
	"order_submission_tool/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// editing order before submission, i.e., quantities of items ordered
func EditOrderHandler(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := db.Collection("orders")
    _, err = collection.UpdateOne(ctx, bson.M{"ID": order.ID}, bson.M{"$set": order})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "the order has been updated successfully"})
}

// deleting an order, i.e., clearing cart
func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := db.Collection("orders")
    _, err = collection.DeleteOne(ctx, bson.M{"ID": order.ID})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "the order has been deleted successfully"})
}

// submitting order to discord channel for processing
func SendOrderHandler(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    orderDetails := ""
    for _, item := range order.Cart {
        // getting each item ID from cart
        itemID, err := primitive.ObjectIDFromHex(item.ItemID.Hex())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // getting item details from its ID
        itemDetails, err := GetItemByID(itemID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        orderDetails += fmt.Sprintf("%s: %d\n", itemDetails.Name, item.Quantity)
    }
    fullMessage := fmt.Sprintf("A new order has been received!\nEmail: %s\nMessage: %s\n\nOrder Details:\n%s", order.Email, order.Message, orderDetails)

    err = bot.SendMessage(fullMessage)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

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

// fetching an Item from the db by its ID
func GetItemByID(itemID primitive.ObjectID) (models.Item, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("items")
    var item models.Item
    err := collection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
    if err != nil {
        return models.Item{}, err
    }

    return item, nil
}