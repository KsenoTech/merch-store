package routes

import (
	"github.com/KsenoTech/merch-store/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler) *mux.Router {
	router := mux.NewRouter()

	// Роуты для аутентификации
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/validate", authHandler.ValidateToken).Methods("GET")

	return router
}
