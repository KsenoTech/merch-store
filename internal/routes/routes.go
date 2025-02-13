package routes

import (
	"net/http"

	"github.com/KsenoTech/merch-store/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler, transactionHandler *handlers.TransactionHandler, jwtMiddleware func(http.Handler) http.Handler) *mux.Router {
	router := mux.NewRouter()

	// Роуты для аутентификации
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/validate", authHandler.ValidateToken).Methods("GET")

	// Защищенные маршруты
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(jwtMiddleware)
	protected.HandleFunc("/transfer", transactionHandler.TransferCoins).Methods("POST")
	protected.HandleFunc("/buy", transactionHandler.BuyMerch).Methods("POST")

	return router
}
