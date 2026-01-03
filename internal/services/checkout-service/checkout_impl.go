package checkoutservice

import (
	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
)

type CheckoutService struct {
	queries       *db.Queries
	logger        *loggermanager.Logger
	rabbitManager rabbitmqmanager.IRabbitMQManager
}

func NewCheckoutService(queries *db.Queries, logger *loggermanager.Logger, rabbitManager rabbitmqmanager.IRabbitMQManager) ICheckoutService {
	return &CheckoutService{
		queries:       queries,
		logger:        logger,
		rabbitManager: rabbitManager,
	}
}

func (cs *CheckoutService) CreatePaymentIntent(amount, currency, reference, callbackURL, merchantID string, rabbitMgr rabbitmqmanager.IRabbitMQManager) error {
	return nil
}
