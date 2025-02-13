package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/KsenoTech/merch-store/internal/middleware"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
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

	// Проверка баланса отправителя
	var balance int
	err := h.DB.QueryRow("SELECT coins FROM users WHERE id = $1", fromUserID).Scan(&balance)
	if err != nil {
		http.Error(w, "Failed to get user balance", http.StatusInternalServerError)
		return
	}
	if balance < transferRequest.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию
	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Списываем монеты у отправителя
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", transferRequest.Amount, fromUserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to deduct coins", http.StatusInternalServerError)
		return
	}

	// Зачисляем монеты получателю
	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", transferRequest.Amount, transferRequest.ToUserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to add coins", http.StatusInternalServerError)
		return
	}

	// Логируем транзакцию
	_, err = tx.Exec("INSERT INTO coin_transfers (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", fromUserID, transferRequest.ToUserID, transferRequest.Amount)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to log transfer", http.StatusInternalServerError)
		return
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coins transferred successfully"))
}

// Покупка мерча
func (h *TransactionHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	var purchaseRequest struct {
		ItemName string `json:"item_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&purchaseRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем цену товара
	var price int
	err := h.DB.QueryRow("SELECT price FROM merch WHERE name = $1", purchaseRequest.ItemName).Scan(&price)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Проверяем баланс пользователя
	var balance int
	err = h.DB.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		http.Error(w, "Failed to get user balance", http.StatusInternalServerError)
		return
	}
	if balance < price {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию
	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Списываем монеты у пользователя
	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to deduct coins", http.StatusInternalServerError)
		return
	}

	// Логируем покупку
	_, err = tx.Exec("INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)", userID, purchaseRequest.ItemName, price)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to log purchase", http.StatusInternalServerError)
		return
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item purchased successfully"))
}
