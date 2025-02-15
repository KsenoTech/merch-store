package migrations

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// Embed migrations file
//
//go:embed up.sql
var migrations embed.FS

// ApplyMigrations применяет SQL-миграции из файла up.sql
func ApplyMigrations(db *gorm.DB) error {

	// Проверяем существование таблиц
	if tableExists(db, "users") && tableExists(db, "merch") && tableExists(db, "coin_transfers") && tableExists(db, "purchases") {
		log.Println("Database is already migrated")
		return nil
	}

	// Читаем содержимое файла up.sql
	migrationContent, err := migrations.ReadFile("up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Разделяем миграции по точке с запятой (если есть несколько запросов)
	statements := strings.Split(string(migrationContent), ";")

	// Выполняем каждый запрос
	for _, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement != "" {
			if err := db.Exec(trimmedStatement).Error; err != nil {
				return fmt.Errorf("failed to execute migration: %w", err)
			}
		}
	}

	log.Println("Migrations applied successfully")
	return nil
}

// Проверка существования таблицы
func tableExists(db *gorm.DB, tableName string) bool {
	var count int
	db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?`, tableName).Scan(&count)
	return count > 0
}
