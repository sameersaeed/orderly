package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"order_submission_tool/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// (admin only) creating items for store
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
    var item models.Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if item.ID == primitive.NilObjectID {  
        item.ID = primitive.NewObjectID()
    }
    
    // saving item to db
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("items")
    _, err = collection.InsertOne(ctx, item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Print("LOG: successfully created item: ", item.ID)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "item created successfully"})
}

// (admin only) editing item in store to change its price / description
func EditItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID string `json:"id"`
        Name string `json:"name"`
        Price float64 `json:"price"`
        Description string `json:"description"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // converting string ID to ObjectID
	itemID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "ERROR: response contains an invalid ID", http.StatusBadRequest)
		return
	}

	collection := db.Collection("items")
	filter := bson.M{"_id": itemID}
    update := bson.M{"$set": bson.M{
        "name": req.Name,
        "price": req.Price,
        "description": req.Description,
    }}

    // editing item in db using its ID
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "ERROR: failed to delete item", http.StatusInternalServerError)
		return
	}
    log.Print("LOG: successfully edited item: ", itemID, result)

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

    // converting string ID to ObjectID
	itemID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "ERROR: response contains an invalid ID", http.StatusBadRequest)
		return
	}

	collection := db.Collection("items")
	filter := bson.M{"_id": itemID}

    // deleting item from db using its ID
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "ERROR: failed to delete item", http.StatusInternalServerError)
		return
	}
    log.Print("LOG: successfully deleted item: ", itemID, result)
    
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item was deleted successfully!"))
}

// fetching items to display on store
func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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