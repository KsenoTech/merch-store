package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

func (r *TransactionRepository) GetInventory(userID int) ([]map[string]int, error) {
	rows, err := r.db.Query(`
        SELECT item_name, COUNT(*) AS quantity
        FROM purchases
        WHERE user_id = $1
        GROUP BY item_name
    `, userID)
	if err != nil {
		return nil, errors.New("failed to fetch inventory")
	}
	defer rows.Close()

	var inventory []map[string]int
	for rows.Next() {
		var itemName string
		var quantity int
		if err := rows.Scan(&itemName, &quantity); err != nil {
			return nil, errors.New("failed to scan inventory row")
		}
		inventory = append(inventory, map[string]int{itemName: quantity})
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("failed to iterate over inventory rows")
	}

	return inventory, nil
}

func (r *TransactionRepository) GetTransactionHistory(userID int) (map[string][]map[string]interface{}, error) {
	receivedRows, err := r.db.Query(`
        SELECT u.username, ct.amount
        FROM coin_transfers ct
        JOIN users u ON ct.from_user_id = u.id
        WHERE ct.to_user_id = $1
    `, userID)
	if err != nil {
		return nil, errors.New("failed to fetch received transactions")
	}
	defer receivedRows.Close()

	sentRows, err := r.db.Query(`
        SELECT u.username, ct.amount
        FROM coin_transfers ct
        JOIN users u ON ct.to_user_id = u.id
        WHERE ct.from_user_id = $1
    `, userID)
	if err != nil {
		return nil, errors.New("failed to fetch sent transactions")
	}
	defer sentRows.Close()

	var received []map[string]interface{}
	var sent []map[string]interface{}

	// Собираем полученные транзакции
	for receivedRows.Next() {
		var username string
		var amount int
		if err := receivedRows.Scan(&username, &amount); err != nil {
			return nil, errors.New("failed to scan received row")
		}
		received = append(received, map[string]interface{}{"fromUser": username, "amount": amount})
	}

	// Проверяем ошибки после обработки полученных транзакций
	if err := receivedRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over received transaction rows: %w", err)
	}

	// Собираем отправленные транзакции
	for sentRows.Next() {
		var username string
		var amount int
		if err := sentRows.Scan(&username, &amount); err != nil {
			return nil, errors.New("failed to scan sent row")
		}
		sent = append(sent, map[string]interface{}{"toUser": username, "amount": amount})
	}

	// Проверяем ошибки после обработки отправленных транзакций
	if err := sentRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over sent transaction rows: %w", err)
	}

	return map[string][]map[string]interface{}{
		"received": received,
		"sent":     sent,
	}, nil
}

// func (r *TransactionRepository) GetTransactionHistory(userID int) (map[string][]map[string]int, error) {
// 	receivedRows, err := r.db.Query(`
//         SELECT from_user_id, amount
//         FROM coin_transfers
//         WHERE to_user_id = $1
//     `, userID)
// 	if err != nil {
// 		return nil, errors.New("failed to fetch received transactions")
// 	}
// 	defer receivedRows.Close()

// 	sentRows, err := r.db.Query(`
//         SELECT to_user_id, amount
//         FROM coin_transfers
//         WHERE from_user_id = $1
//     `, userID)
// 	if err != nil {
// 		return nil, errors.New("failed to fetch sent transactions")
// 	}
// 	defer sentRows.Close()

// 	var received []map[string]int
// 	var sent []map[string]int

// 	// Собираем полученные транзакции
// 	for receivedRows.Next() {
// 		var fromUserID int
// 		var amount int
// 		if err := receivedRows.Scan(&fromUserID, &amount); err != nil {
// 			return nil, errors.New("failed to scan received row")
// 		}
// 		received = append(received, map[string]int{"fromUser": fromUserID, "amount": amount})
// 	}

// 	// Проверяем ошибки после обработки полученных транзакций
// 	if err := receivedRows.Err(); err != nil {
// 		return nil, errors.New("failed to iterate over received transaction rows")
// 	}

// 	// Собираем отправленные транзакции
// 	for sentRows.Next() {
// 		var toUserID int
// 		var amount int
// 		if err := sentRows.Scan(&toUserID, &amount); err != nil {
// 			return nil, errors.New("failed to scan sent row")
// 		}
// 		sent = append(sent, map[string]int{"toUser": toUserID, "amount": amount})
// 	}

// 	// Проверяем ошибки после обработки отправленных транзакций
// 	if err := sentRows.Err(); err != nil {
// 		return nil, errors.New("failed to iterate over sent transaction rows")
// 	}

// 	return map[string][]map[string]int{
// 		"received": received,
// 		"sent":     sent,
// 	}, nil
// }
