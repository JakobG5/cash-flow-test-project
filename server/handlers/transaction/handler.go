package transaction

import (
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"
	transactionservice "cash-flow-financial/internal/services/transaction-service"
)

type TransactionHandler struct {
	TRANSACTIONSERVICE transactionservice.ITransactionService
	config             *models.Config
	logger             *loggermanager.Logger
}

func NewTransactionHandler(transactionservice transactionservice.ITransactionService, config *models.Config, logger *loggermanager.Logger) *TransactionHandler {
	return &TransactionHandler{
		TRANSACTIONSERVICE: transactionservice,
		config:             config,
		logger:             logger,
	}
}
