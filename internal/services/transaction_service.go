package service

import "github.com/KsenoTech/merch-store/internal/repository"

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, userRepo: userRepo}
}
