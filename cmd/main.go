package main

import (
	"log"
	"net/http"

	"github.com/KsenoTech/merch-store/internal/database"
	"github.com/KsenoTech/merch-store/internal/handlers"
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

	// Подключаемся к базе данных
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Создаем репозиторий пользователей
	userRepo := repository.NewUserRepository(db)

	// Создаем сервис аутентификации
	authService := service.NewAuthService("your_secret_key", userRepo)

	// Создаем обработчик аутентификации
	authHandler := handlers.NewAuthHandler(authService)

	// Настраиваем роуты
	router := routes.SetupRoutes(authHandler)

	// Запускаем сервер
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
