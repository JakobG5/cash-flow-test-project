package account

import (
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"
	accountservice "cash-flow-financial/internal/services/account-service"
)

type AccountHandler struct {
	accountService accountservice.IAccountService
	config         *models.Config
	logger         *loggermanager.Logger
}

func NewAccountHandler(accountService accountservice.IAccountService, config *models.Config, logger *loggermanager.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		config:         config,
		logger:         logger,
	}
}
