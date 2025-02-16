package routes

import (
	"net/http"

	"github.com/KsenoTech/merch-store/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler, transactionHandler *handlers.TransactionHandler, jwtMiddleware func(http.Handler) http.Handler) http.Handler {
	router := mux.NewRouter()

	// Auth endpoint
	router.HandleFunc("/api/auth", authHandler.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(jwtMiddleware)

	// Получение информации о пользователе
	protected.HandleFunc("/info", transactionHandler.GetUserInfo).Methods("GET")

	// Отправка монеток
	protected.HandleFunc("/sendCoin", transactionHandler.SendCoins).Methods("POST")

	// Покупка мерча
	protected.HandleFunc("/buy/{item}", transactionHandler.BuyMerch).Methods("GET")

	return router
}
