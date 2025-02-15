package repository

import (
	"database/sql"
	"errors"
	"log"
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
	log.Printf("Getting item price for itemName: %s", itemName)

	var price int
	err := tx.QueryRow("SELECT price FROM merch WHERE name = $1", itemName).Scan(&price)
	if err != nil {
		log.Printf("Error getting item price for itemName: %s. Error: %v", itemName, err)
		return 0, errors.New("item not found")
	}
	log.Printf("Item price for %s is %d coins", itemName, price)

	return price, nil
}

// Логирование покупки
func (r *TransactionRepository) LogPurchase(tx *sql.Tx, userID int, itemName string, price int) error {
	_, err := tx.Exec("INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)", userID, itemName, price)
	if err != nil {
		return errors.New("failed to log purchase")
	}
	return nil
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
	log.Printf("Fetching transaction history for userID: %d", userID)

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
		log.Printf("Error fetching transaction history for userID: %d. Error: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var transactionHistory []string
	for rows.Next() {
		var description string
		if err := rows.Scan(&description); err != nil {
			log.Printf("Error scanning transaction history row for userID: %d. Error: %v", userID, err)
			return nil, err
		}
		transactionHistory = append(transactionHistory, description)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over transaction history rows for userID: %d. Error: %v", userID, err)
		return nil, err
	}

	log.Printf("Error iterating over transaction history rows for userID: %d. Error: %v", userID, err)
	return transactionHistory, nil
}
