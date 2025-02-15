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

	"github.com/KsenoTech/merch-store/internal/migrations"
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

	// Применяем миграции
	if err := migrations.ApplyMigrations(db); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied or already up-to-date")

	sqlDB, err := db.DB() // Получаем *sql.DB и проверяем ошибку
	if err != nil {
		log.Fatalf("Failed to get underlying SQL DB: %v", err)
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(sqlDB)

	// Создаем сервисы
	authService := service.NewAuthService("SECRET_KEY", userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Создаем middleware
	jwtMiddleware := middleware.JWTMiddleware("SECRET_KEY") // Вызываем JWTMiddleware с secretKey

	// Настраиваем роуты
	router := routes.SetupRoutes(authHandler, transactionHandler, jwtMiddleware)

	// Запускаем сервер
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
