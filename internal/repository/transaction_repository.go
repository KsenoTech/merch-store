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

func (r *TransactionRepository) BeginTx() (*sql.Tx, error) {
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

// Получение истории покупок пользователя
func (r *TransactionRepository) GetPurchasedItems(tx *sql.Tx, userID int) ([]string, error) {
	rows, err := tx.Query("SELECT item_name FROM purchases WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchasedItems []string
	for rows.Next() {
		var itemName string
		if err := rows.Scan(&itemName); err != nil {
			return nil, err
		}
		purchasedItems = append(purchasedItems, itemName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return purchasedItems, nil
}

// Получение истории транзакций пользователя
func (r *TransactionRepository) GetTransactionHistory(tx *sql.Tx, userID int) ([]string, error) {
	rows, err := tx.Query(`
        SELECT 
            CASE 
                WHEN from_user_id = $1 THEN 'Sent to ' || to_user_id || ': ' || amount 
                WHEN to_user_id = $1 THEN 'Received from ' || from_user_id || ': ' || amount 
            END AS transaction_description
        FROM coin_transfers
        WHERE from_user_id = $1 OR to_user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionHistory []string
	for rows.Next() {
		var description string
		if err := rows.Scan(&description); err != nil {
			return nil, err
		}
		transactionHistory = append(transactionHistory, description)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactionHistory, nil
}
