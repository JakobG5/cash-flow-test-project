package checkout

import (
	loggermanager "cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	"cash-flow-financial/internal/models"
	accountservice "cash-flow-financial/internal/services/account-service"
	checkoutservice "cash-flow-financial/internal/services/checkout-service"
)

type CheckoutHandler struct {
	checkoutService checkoutservice.ICheckoutService
	accountService  accountservice.IAccountService
	rabbitManager   rabbitmqmanager.IRabbitMQManager
	config          *models.Config
	logger          *loggermanager.Logger
}

func NewCheckoutHandler(checkoutService checkoutservice.ICheckoutService, accountService accountservice.IAccountService, config *models.Config, logger *loggermanager.Logger, rabbitMgr rabbitmqmanager.IRabbitMQManager) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
		accountService:  accountService,
		rabbitManager:   rabbitMgr,
		config:          config,
		logger:          logger,
	}
}
