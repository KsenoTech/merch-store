package service

import (
	"context"
	"errors"

	"github.com/KsenoTech/merch-store/internal/models"
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
	tx, err := s.transactionRepo.Begin()
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

// Получение списка купленных товаров
func (s *TransactionService) GetPurchasedItems(ctx context.Context, userID int) ([]models.PurchasedItem, error) {
	return s.transactionRepo.GetPurchasedItems(ctx, userID)
}

// Получение истории транзакций (отправленных и полученных монет)
func (s *TransactionService) GetTransactionHistory(ctx context.Context, userID int) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionHistory(ctx, userID)
}

// Покупка мерча
func (s *TransactionService) BuyMerch(userID int, itemName string) error {
	// Начинаем транзакцию
	tx, err := s.transactionRepo.Begin()
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

	// Получаем цену товара
	price, err := s.transactionRepo.GetItemPrice(tx, itemName)
	if err != nil {
		return errors.New("item not found")
	}

	// Проверяем баланс пользователя
	balance, err := s.userRepo.GetUserBalance(tx, userID)
	if err != nil {
		return errors.New("failed to get user balance")
	}
	if balance < price {
		return errors.New("insufficient funds")
	}

	// Списываем монеты у пользователя
	err = s.userRepo.DeductCoins(tx, userID, price)
	if err != nil {
		return errors.New("failed to deduct coins")
	}

	// Логируем покупку
	err = s.transactionRepo.LogPurchase(tx, userID, itemName, price)
	if err != nil {
		return errors.New("failed to log purchase")
	}

	return nil
}
