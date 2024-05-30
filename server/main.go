package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"order_submission_tool/bot"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func main() {
	err := godotenv.Load("../client/.env")
	if err != nil {
		log.Fatal("ERROR: could not fetch env vars from .env file - please make sure the path is correct")
	}

	log.Printf("server is now running on: %s:%s\n", os.Getenv("HOST_URL"), os.Getenv("HOST_PORT"))

	err = bot.Start()
	if err != nil {
		log.Fatalf("ERROR: could not start bot: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/send", sendHandler).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":3001", handler))
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	var data FormData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("A new order has been received!\nName: %s\nEmail: %s\nOrder details: %s", data.Name, data.Email, data.Message)

	err = bot.SendMessage(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order was sent to Discord channel" + os.Getenv("DISCORD_CHANNEL_ID") + "successfully!")
}
