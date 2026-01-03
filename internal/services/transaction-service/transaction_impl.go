package transactionservice

import (
	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
)

type TransactionService struct {
	queries *db.Queries
	logger  *loggermanager.Logger
}

func NewTransactionService(queries *db.Queries, logger *loggermanager.Logger) ITransactionService {
	return &TransactionService{
		queries: queries,
		logger:  logger,
	}
}

func (ts *TransactionService) GetPaymentStatus(transactionID string) error {
	return nil
}

func (ts *TransactionService) ProcessPayment(transactionID string) error {
	return nil
}
