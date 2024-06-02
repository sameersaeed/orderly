package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"order_submission_tool/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    err := godotenv.Load("../client/.env")
    if err != nil {
        log.Fatal("ERROR: could not fetch env vars from .env file - please make sure the path is correct")
    }

    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("ERROR: could not connect to MongoDB: %v", err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("ERROR: could not connect to MongoDB: %v", err)
    }

    db := client.Database("shop")
    handlers.SetDatabase(db)

    log.Printf("server is now running on: %s:%s\n", os.Getenv("HOST_URL"), os.Getenv("SERVER_PORT"))

    router := mux.NewRouter()
    router.HandleFunc("/jwt", handlers.SessionHandler).Methods("GET")

    router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
    router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

    router.HandleFunc("/getOrders", handlers.GetOrdersHandler).Methods("GET")
	router.HandleFunc("/editOrder", handlers.EditOrderHandler).Methods("POST")
    router.HandleFunc("/deleteOrder", handlers.DeleteOrderHandler).Methods("DELETE")

    router.HandleFunc("/getItems", handlers.GetItemsHandler).Methods("GET")
    router.HandleFunc("/createItem", handlers.CreateItemHandler).Methods("POST")
    router.HandleFunc("/editItem", handlers.EditItemHandler).Methods("POST")
    router.HandleFunc("/deleteItem", handlers.DeleteItemHandler).Methods("DELETE")
    
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{ os.Getenv("HOST_URL") + ":" + os.Getenv("HOST_PORT") },
        AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    handler := c.Handler(router)

    log.Fatal(http.ListenAndServe(":3001", handler))
}
