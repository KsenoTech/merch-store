package repository

import (
	"database/sql"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/KsenoTech/merch-store/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	log.Printf("Creating new user: %+v", user)
	result := r.db.Create(user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return errors.New("failed to create user")
	}
	log.Printf("User created successfully: %+v", user)
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	log.Printf("Fetching user by username: %s", username)
	var user models.User
	err := r.db.Model(&models.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Пользователь не найден
		}
		return nil, errors.New("failed to get user by username")
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// Получение баланса пользователя
func (r *UserRepository) GetUserBalance(tx *sql.Tx, userID int) (int, error) {
	log.Printf("Getting user balance for userID: %d", userID)

	var balance int
	var query string

	// Если транзакция не передана, используем основное соединение
	if tx == nil {
		// Получаем *sql.DB из GORM
		sqlDB, err := r.db.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
			return 0, errors.New("failed to get underlying SQL DB")
		}

		query = "SELECT coins FROM users WHERE id = $1"
		err = sqlDB.QueryRow(query, userID).Scan(&balance)

		if err != nil {
			log.Printf("Error getting user balance for userID: %d. Error: %v", userID, err)
			return 0, errors.New("failed to get user balance")
		}
	} else {
		query = "SELECT coins FROM users WHERE id = $1"
		err := tx.QueryRow(query, userID).Scan(&balance)
		if err != nil {
			log.Printf("Error getting user balance for userID: %d in transaction. Error: %v", userID, err)
			return 0, errors.New("failed to get user balance")
		}
	}

	log.Printf("User balance for userID: %d is %d coins", userID, balance)
	return balance, nil
}

// Списание монет
func (r *UserRepository) DeductCoins(tx *sql.Tx, userID, amount int) error {
	_, err := tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, userID)
	if err != nil {
		return errors.New("failed to deduct coins")
	}
	return nil
}

func (r *UserRepository) AddCoins(tx *sql.Tx, userID, amount int) error {
	_, err := tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, userID)
	return err
}
