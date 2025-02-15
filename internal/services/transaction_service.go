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

// Получение истории покупок пользователя
func (s *TransactionService) GetPurchasedItems(userID int) ([]string, error) {
	tx, err := s.transactionRepo.BeginTx()
	if err != nil {
		return nil, errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	return s.transactionRepo.GetPurchasedItems(tx, userID)
}

// Получение истории транзакций пользователя
func (s *TransactionService) GetTransactionHistory(userID int) ([]string, error) {
	log.Printf("Getting transaction history for userID: %d", userID)

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

	history, err := s.transactionRepo.GetTransactionHistory(tx, userID)
	if err != nil {
		log.Printf("Error getting transaction history for userID: %d. Error: %v", userID, err)
		return nil, errors.New("failed to get transaction history")
	}

	log.Printf("Transaction history for userID: %d is %v", userID, history)
	return history, nil
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
