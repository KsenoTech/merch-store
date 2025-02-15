package repository

import (
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *TransactionRepository) LogTransfer(tx *sql.Tx, fromUserID, toUserID, amount int) error {
	_, err := tx.Exec("INSERT INTO coin_transfers (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", fromUserID, toUserID, amount)
	return err
}

func (r *TransactionRepository) GetItemPrice(tx *sql.Tx, itemName string) (int, error) {
	var price int
	err := tx.QueryRow("SELECT price FROM merch WHERE name = $1", itemName).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (r *TransactionRepository) LogPurchase(tx *sql.Tx, userID int, itemName string, price int) error {
	_, err := tx.Exec("INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)", userID, itemName, price)
	return err
}
