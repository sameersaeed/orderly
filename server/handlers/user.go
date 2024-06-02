package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order_submission_tool/models"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

// handling registration requests to add new users to db
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // saving user info to db
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
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

// handling user login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    collection := db.Collection("users")
    err = collection.FindOne(ctx, bson.M{"email":  strings.ToLower(user.Email)}).Decode(&user)
    if err != nil {
        http.Error(w, "ERROR: invalid email or password", http.StatusUnauthorized)
        return
    }

    // creating a JWT token for the user, using their email
    token, expirationTime, jwtErr := generateJWT(user.Email)
    if jwtErr != nil {
        http.Error(w, "ERROR: could not generate JWT token", http.StatusInternalServerError)
        log.Printf("LOG: JWT token generation error for user %s: %s\n", user.Email, jwtErr)
        return
    }

    user.Token = token

    // adding JWT token as a cookie on the user's browser
    http.SetCookie(w, &http.Cookie{
        Name:    "jwtToken",
        Value:   token,
        Expires: expirationTime,
    })

    w.Header().Set("Authorization", "Bearer " + token)
    w.WriteHeader(http.StatusOK)
    log.Printf("LOG: user %s logged in successfully\n", user.Email)

    json.NewEncoder(w).Encode(user)
}

// validating whether user has authorization to access the shop
func SessionHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    tokenString := r.Header.Get("Authorization")

    if tokenString == "" {
        w.WriteHeader(http.StatusUnauthorized)
        http.Error(w, "ERROR: no JWT token was found", http.StatusUnauthorized)
        return
    }
    tokenString = tokenString[len("Bearer "):]

    // verifying JWT token
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_KEY")), nil
    })
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    if !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprint(w, "ERROR: invalid JWT token was provided")
        return
    }

    // validating token hasn't expired yet
    expirationTime := time.Unix(int64(claims["expiration"].(float64)), 0)
    if time.Now().After(expirationTime) {
        // if token has expired, logout the user
        w.WriteHeader(http.StatusUnauthorized)
        http.SetCookie(w, &http.Cookie{
            Name:    "jwtToken",
            Value:   "",
            Expires: time.Now().Add(-time.Hour),
        })
        return
    }

	response := map[string]interface{}{
		"token":     tokenString,
		"expiresAt": expirationTime,
		"email":     claims["email"],
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "ERROR: could not marshal JSON response: %v", err)
		return
	}
    
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)

}

// generating JWT token for the user using their email and a secret key
func generateJWT(email string) (string, time.Time, error) {
    // generates a new token which expires in 24hrs, with user's email as a claim
    expiration := time.Now().Add(time.Hour * 24)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ 
        "email": email, 
        "expiration": expiration.Unix(), 
    })

    // signing token with JWT secret
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
    if err != nil {
        return "", time.Time{}, err
    }

    return tokenString, expiration, nil
}