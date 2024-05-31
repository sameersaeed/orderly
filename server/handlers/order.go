package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order_submission_tool/bot"
	"order_submission_tool/models"
)

type Order struct {
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Message string        `json:"message"`
	Cart    []models.Item `json:"cart"`
}

func SendOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderDetails := ""
	for _, item := range order.Cart {
		orderDetails += fmt.Sprintf("%s: %.2f\n", item.Name, item.Price)
	}
	fullMessage := fmt.Sprintf("A new order has been received!\nName: %s\nEmail: %s\nOrder details: %s\n\nOrder Details:\n%s", order.Name, order.Email, order.Message, orderDetails)

	err = bot.SendMessage(fullMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order sent to Discord channel successfully"})
}