package repository

import (
	"database/sql"

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

func (r *UserRepository) GetUserBalance(tx *sql.Tx, userID int) (int, error) {
	var balance int
	err := tx.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *UserRepository) DeductCoins(tx *sql.Tx, userID, amount int) error {
	_, err := tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, userID)
	return err
}

func (r *UserRepository) AddCoins(tx *sql.Tx, userID, amount int) error {
	_, err := tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, userID)
	return err
}
