package checkout

import (
	loggermanager "cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	"cash-flow-financial/internal/models"
	checkoutservice "cash-flow-financial/internal/services/checkout-service"
)

type CheckoutHandler struct {
	checkoutService checkoutservice.ICheckoutService
	rabbitManager   rabbitmqmanager.IRabbitMQManager
	config          *models.Config
	logger          *loggermanager.Logger
}

func NewCheckoutHandler(checkoutService checkoutservice.ICheckoutService, config *models.Config, logger *loggermanager.Logger, rabbitMgr rabbitmqmanager.IRabbitMQManager) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
		rabbitManager:   rabbitMgr,
		config:          config,
		logger:          logger,
	}
}
