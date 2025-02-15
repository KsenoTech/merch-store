package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KsenoTech/merch-store/internal/middleware"
	services "github.com/KsenoTech/merch-store/internal/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Перевод монет между пользователями
func (h *TransactionHandler) TransferCoins(w http.ResponseWriter, r *http.Request) {
	var transferRequest struct {
		ToUserID int `json:"to_user_id"`
		Amount   int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&transferRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	fromUserID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.TransferCoins(fromUserID, transferRequest.ToUserID, transferRequest.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coins transferred successfully"))
}

// Покупка мерча
func (h *TransactionHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var purchaseRequest struct {
		ItemName string `json:"item_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&purchaseRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Processing BuyMerch request for userID: %d, itemName: %s", userID, purchaseRequest.ItemName)

	err := h.service.BuyMerch(userID, purchaseRequest.ItemName)
	if err != nil {
		log.Printf("Error processing BuyMerch request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("BuyMerch request processed successfully for userID: %d, itemName: %s", userID, purchaseRequest.ItemName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item purchased successfully"))
}

// Получение списка купленных товаров
func (h *TransactionHandler) GetPurchasedItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.service.GetPurchasedItems(userID)
	if err != nil {
		http.Error(w, "Failed to get purchased items", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

// Получение истории транзакций
func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Processing GetTransactionHistory request for userID: %d", userID)

	history, err := h.service.GetTransactionHistory(userID)
	if err != nil {
		log.Printf("Error processing GetTransactionHistory request for userID: %d. Error: %v", userID, err)
		http.Error(w, "Failed to get transaction history", http.StatusInternalServerError)
		return
	}
	log.Printf("GetTransactionHistory request processed successfully for userID: %d. History: %v", userID, history)

	json.NewEncoder(w).Encode(history)
}
