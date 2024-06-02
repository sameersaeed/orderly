package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"order_submission_tool/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// (admin only) creating items for store
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
    var item models.Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("items")
    _, err = collection.InsertOne(ctx, item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "item created successfully"})
}

// (admin only) editing item in store to change its price / description
func EditItemHandler(w http.ResponseWriter, r *http.Request) {
    var item models.Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := db.Collection("items")
    _, err = collection.UpdateOne(ctx, bson.M{"ID": item.ID}, bson.M{"$set": item})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "item updated successfully"})
}

// (admin only) deleting item from store
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID string `json:"id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := db.Collection("items")
	filter := bson.M{"_id": req.ID}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "ERROR: failed to delete item", http.StatusInternalServerError)
		return
	}

    log.Print("deleted item: ", result);

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The item was deleted successfully!"))
}

// fetching items to display on store
func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("items")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)
    var items []models.Item
    for cursor.Next(ctx) {
        var item models.Item
        err := cursor.Decode(&item)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        items = append(items, item)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(items)
}
