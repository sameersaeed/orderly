package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"order_submission_tool/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("users")

	// validating whether email already exists in db
	existingUser := models.User{}
	err = collection.FindOne(ctx, bson.M{"email": strings.ToLower(user.Email)}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "ERROR: could not register user. provided email is already in use", http.StatusBadRequest)
		return
	}

	// if not, can register new user
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	log.Printf("LOG: a new user has registered (ID: %s)\n", result.InsertedID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("users")
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "ERROR: invalid email or password", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("LOG: user %s logged in successfully\n", user.Email)

	json.NewEncoder(w).Encode(user)
}
