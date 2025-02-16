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

func ApplyMigrations(db *gorm.DB) error {
	// Проверяем, существует ли таблица schema_migrations
	var count int64
	db.Model(&struct{}{}).Table("information_schema.tables").
		Where("table_name = ?", "schema_migrations").Count(&count)

	if count == 0 {
		log.Println("Applying migrations...")

		// Читаем содержимое файла up.sql
		migrationContent, err := migrations.ReadFile("up.sql")
		if err != nil {
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		// Разделяем миграции по точке с запятой
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

		// Создаем таблицу schema_migrations после успешного выполнения миграций
		if err := db.Exec(`CREATE TABLE schema_migrations (version TEXT PRIMARY KEY)`).Error; err != nil {
			return fmt.Errorf("failed to create migration table: %w", err)
		}

		if err := db.Exec(`INSERT INTO schema_migrations (version) VALUES ('v1')`).Error; err != nil {
			return fmt.Errorf("failed to insert migration version: %w", err)
		}

		log.Println("Migrations applied successfully")
	} else {
		log.Println("Migrations already applied")
	}

	return nil
}

// Проверка существования таблицы schema_migrations
func migrationTableExists(db *gorm.DB) bool {
	var count int64
	db.Model(&struct{}{}).Table("information_schema.tables").
		Where("table_name = ?", "schema_migrations").Count(&count)
	return count > 0
}

// Создание таблицы schema_migrations
func createMigrationTable(db *gorm.DB) error {
	return db.Exec(`
        CREATE TABLE schema_migrations (
            version TEXT PRIMARY KEY
        );
        INSERT INTO schema_migrations (version) VALUES ('v1');
    `).Error
}

// Выполнение миграций
func executeMigrations(db *gorm.DB, sql string) error {
	// Разделяем SQL-запросы по точке с запятой
	statements := splitSQL(sql)

	// Выполняем каждый запрос
	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute migration statement: %w", err)
		}
	}
	return nil
}

// Разделение SQL-запросов
func splitSQL(sql string) []string {
	var statements []string
	var buffer strings.Builder
	inString := false

	for _, char := range sql {
		if char == '\'' || char == '"' {
			inString = !inString
		}
		if char == ';' && !inString {
			statements = append(statements, buffer.String())
			buffer.Reset()
		} else {
			buffer.WriteRune(char)
		}
	}
	if buffer.Len() > 0 {
		statements = append(statements, buffer.String())
	}
	return statements
}
