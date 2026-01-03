package handlers

import (
	"cash-flow-financial/internal/models"
	accountservice "cash-flow-financial/internal/services/account-service"
	checkoutservice "cash-flow-financial/internal/services/checkout-service"
	transactionservice "cash-flow-financial/internal/services/transaction-service"

	"github.com/docker/docker/daemon/logger"
)

type Handlers struct {
	CHECKOUTSERVICE    checkoutservice.ICheckoutService
	ACCOUNTSERVICE     accountservice.IAccountService
	TRANSACTIONSERVICE transactionservice.ITransactionService
	CONFIG             models.Config
	LOGGER             logger.Logger
}

func NewHandlers(checkoutservice checkoutservice.ICheckoutService, accountservice accountservice.IAccountService, transactionservice transactionservice.ITransactionService, config models.Config, logger logger.Logger) *Handlers {
	return &Handlers{
		CHECKOUTSERVICE:    checkoutservice,
		ACCOUNTSERVICE:     accountservice,
		TRANSACTIONSERVICE: transactionservice,
		CONFIG:             config,
		LOGGER:             logger,
	}
}
