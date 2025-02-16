package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func ConnectDB() (*gorm.DB, error) {
// 	// Получаем строку подключения
// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		return nil, fmt.Errorf("DATABASE_URL is not set")
// 	}

// 	// Разделяем строку подключения на части
// 	dsnParts := parseDSN(dbURL)
// 	if dsnParts == nil {
// 		return nil, fmt.Errorf("invalid DATABASE_URL")
// 	}

// 	// Проверяем существование базы данных и создаём её, если она не существует
// 	if err := ensureDatabaseExists(dsnParts); err != nil {
// 		return nil, err
// 	}

// 	// Подключаемся к базе данных
// 	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

//		log.Println("Connected to the database successfully")
//		return db, nil
//	}

func ConnectDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем сырое SQL-соединение
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Создание базы данных, если она не существует
func ensureDatabaseExists(dsnParts map[string]string) error {
	// Создаём DSN без имени базы данных
	dsnWithoutDB := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s",
		dsnParts["host"], dsnParts["user"], dsnParts["password"], dsnParts["port"], dsnParts["sslmode"])

	// Подключаемся к PostgreSQL без базы данных
	db, err := sql.Open("postgres", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	// Проверяем существование базы данных
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pg_database WHERE datname = $1", dsnParts["dbname"]).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if count == 0 {
		// Создаём базу данных, если её нет
		log.Printf("Database %s does not exist, creating...", dsnParts["dbname"])
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dsnParts["dbname"]))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s created successfully", dsnParts["dbname"])
	}

	return nil
}

// Парсинг строки подключения
func parseDSN(dsn string) map[string]string {
	parts := make(map[string]string)

	// Удаляем префикс "postgres://"
	dsn = strings.TrimPrefix(dsn, "postgres://")

	// Разделяем строку на userInfo и остальную часть
	if idx := strings.Index(dsn, "@"); idx != -1 {
		userInfo := dsn[:idx]
		dsn = dsn[idx+1:]

		// Разбиваем userInfo на пользователя и пароль
		if split := strings.Split(userInfo, ":"); len(split) == 2 {
			parts["user"] = split[0]
			parts["password"] = split[1]
		}
	}

	// Разделяем хост и параметры
	if idx := strings.Index(dsn, "/"); idx != -1 {
		hostPort := dsn[:idx]
		dsn = dsn[idx+1:]

		// Разбиваем hostPort на хост и порт
		if split := strings.Split(hostPort, ":"); len(split) == 2 {
			parts["host"] = split[0]
			parts["port"] = split[1]
		} else {
			parts["host"] = hostPort
			parts["port"] = "5432" // По умолчанию порт PostgreSQL
		}
	}

	// Извлекаем имя базы данных
	if idx := strings.Index(dsn, "?"); idx != -1 {
		parts["dbname"] = dsn[:idx]
		query := dsn[idx+1:]
		params, _ := url.ParseQuery(query)
		for key, value := range params {
			parts[key] = value[0]
		}
	} else {
		parts["dbname"] = dsn
	}

	// Устанавливаем sslmode по умолчанию, если он не указан
	if _, ok := parts["sslmode"]; !ok {
		parts["sslmode"] = "disable"
	}

	return parts
}
