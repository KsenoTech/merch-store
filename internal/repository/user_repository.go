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
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// Получение баланса пользователя
func (r *UserRepository) GetUserBalance(tx *sql.Tx, userID int) (int, error) {
	log.Printf("Getting user balance for userID: %d", userID)

	var coins int
	err := tx.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&coins)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error UserRepository getting user balance for userID: %d. Error: %v", userID, err)
		return 0, errors.New("failed to get user balance")
	}
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found")
	}

	log.Printf("User balance for userID: %d is %d coins", userID, coins)
	return coins, nil
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
