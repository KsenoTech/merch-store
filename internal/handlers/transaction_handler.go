package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KsenoTech/merch-store/internal/middleware"
	services "github.com/KsenoTech/merch-store/internal/services"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Покупка мерча
func (h *TransactionHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	// Извлекаем userID из контекста
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Извлекаем название товара из параметров URL
	vars := mux.Vars(r)
	itemName := vars["item"]
	if itemName == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Processing BuyMerch request for userID: %d, itemName: %s", userID, itemName)

	// Вызываем сервис для покупки мерча
	err := h.service.BuyMerch(userID, itemName)
	if err != nil {
		log.Printf("Error processing BuyMerch request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item purchased successfully"))
}

func (h *TransactionHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Processing GetUserInfo request for userID: %d", userID)

	info, err := h.service.GetUserInfo(userID)
	if err != nil {
		log.Printf("Error processing GetUserInfo request for userID: %d. Error: %v", userID, err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	log.Printf("GetUserInfo request processed successfully for userID: %d. Info: %+v", userID, info)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (h *TransactionHandler) SendCoins(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	toUserID, err := h.service.GetUserByUsername(req.ToUser)
	if err != nil {
		http.Error(w, "Recipient not found", http.StatusBadRequest)
		return
	}

	err = h.service.TransferCoins(userID, toUserID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coins sent successfully"))
}
