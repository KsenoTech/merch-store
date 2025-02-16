package service

import (
	"errors"
	"log"

	"github.com/KsenoTech/merch-store/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, userRepo: userRepo}
}

func (s *TransactionService) GetUserByUsername(username string) (int, error) {
	log.Printf("Getting user ID for username: %s", username)
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user ID for username: %s. Error: %v", username, err)
		return 0, errors.New("user not found")
	}
	userID := int(user.ID)
	log.Printf("User ID for username %s is %d", username, userID)
	return userID, nil
}

// Перевод монет между пользователями
func (s *TransactionService) TransferCoins(fromUserID, toUserID, amount int) error {
	// Начинаем транзакцию
	tx, err := s.transactionRepo.BeginTx()
	if err != nil {
		return errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Проверяем баланс отправителя
	balance, err := s.userRepo.GetUserBalance(tx, fromUserID)
	if err != nil {
		return errors.New("failed to get user balance")
	}
	if balance < amount {
		return errors.New("insufficient funds")
	}

	// Списываем монеты у отправителя
	err = s.userRepo.DeductCoins(tx, fromUserID, amount)
	if err != nil {
		return errors.New("failed to deduct coins")
	}

	// Зачисляем монеты получателю
	err = s.userRepo.AddCoins(tx, toUserID, amount)
	if err != nil {
		return errors.New("failed to add coins")
	}

	// Логируем транзакцию
	err = s.transactionRepo.LogTransfer(tx, fromUserID, toUserID, amount)
	if err != nil {
		return errors.New("failed to log transfer")
	}

	return nil
}

// Покупка мерча
func (s *TransactionService) BuyMerch(userID int, itemName string) error {

	log.Printf("Starting BuyMerch for userID: %d, itemName: %s", userID, itemName)

	// Начинаем транзакцию
	tx, err := s.transactionRepo.BeginTx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error")
			_ = tx.Rollback()
		} else {
			log.Println("Committing transaction")
			_ = tx.Commit()
		}
	}()

	// Получаем цену товара
	price, err := s.transactionRepo.GetItemPrice(tx, itemName)
	if err != nil {
		log.Printf("Error getting item price for item: %s. Error: %v", itemName, err)
		return errors.New("item not found")
	}
	log.Printf("Item price for %s is %d coins", itemName, price)

	// Проверяем баланс пользователя
	balance, err := s.userRepo.GetUserBalance(tx, userID)
	if err != nil {
		log.Printf("Error getting user balance for userID: %d. Error: %v", userID, err)
		return errors.New("failed to get user balance")
	}
	log.Printf("User balance for userID: %d is %d coins", userID, balance)

	if balance < price {
		log.Printf("Insufficient funds for userID: %d. Required: %d, Available: %d", userID, price, balance)
		return errors.New("insufficient funds")
	}

	// Списываем монеты у пользователя
	err = s.userRepo.DeductCoins(tx, userID, price)
	if err != nil {
		log.Printf("Error deducting coins from userID: %d. Error: %v", userID, err)
		return errors.New("failed to deduct coins")
	}
	log.Printf("Deducted %d coins from userID: %d", price, userID)

	// Логируем покупку
	err = s.transactionRepo.LogPurchase(tx, userID, itemName, price)
	if err != nil {
		log.Printf("Error logging purchase for userID: %d, item: %s. Error: %v", userID, itemName, err)
		return errors.New("failed to log purchase")
	}
	log.Printf("Logged purchase of %s for userID: %d", itemName, userID)

	return nil
}

func (s *TransactionService) GetUserInfo(userID int) (map[string]interface{}, error) {
	log.Printf("Getting user info for userID: %d", userID)

	// Начинаем транзакцию
	tx, err := s.transactionRepo.BeginTx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Получаем баланс пользователя
	balance, err := s.userRepo.GetUserBalance(nil, userID)
	if err != nil {
		log.Printf("Error getting user balance for userID: %d. Error: %v", userID, err)
		return nil, errors.New("failed to get user balance")
	}

	// Получаем инвентарь пользователя
	inventory, err := s.transactionRepo.GetInventory(userID)
	if err != nil {
		log.Printf("Error getting user inventory for userID: %d. Error: %v", userID, err)
		return nil, errors.New("failed to get user inventory")
	}

	// Получаем историю транзакций
	coinHistory, err := s.transactionRepo.GetTransactionHistory(userID)
	if err != nil {
		log.Printf("Error getting transaction history for userID: %d. Error: %v", userID, err)
		return nil, errors.New("failed to get transaction history")
	}

	// Формируем ответ
	response := map[string]interface{}{
		"coins":       balance,
		"inventory":   inventory,
		"coinHistory": coinHistory,
	}

	log.Printf("User info for userID: %d is %+v", userID, response)
	return response, nil
}
