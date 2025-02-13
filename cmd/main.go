package main

import (
	"log"
	"net/http"

	"github.com/KsenoTech/merch-store/internal/database"
	"github.com/KsenoTech/merch-store/internal/handlers"
	"github.com/KsenoTech/merch-store/internal/middleware"
	"github.com/KsenoTech/merch-store/internal/repository"
	"github.com/KsenoTech/merch-store/internal/routes"
	service "github.com/KsenoTech/merch-store/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Подключаемся к базе данных через GORM
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Создаем сервисы
	authService := service.NewAuthService("your_secret_key", userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Настраиваем роуты
	router := routes.SetupRoutes(authHandler, transactionHandler, middleware.JWTMiddleware)

	// Запускаем сервер
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
