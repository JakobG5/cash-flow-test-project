package checkoutservice

import "cash-flow-financial/internal/managers/rabbitmqmanager"

type ICheckoutService interface {
	CreatePaymentIntent(amount, currency, reference, callbackURL, merchantID string, rabbitMgr rabbitmqmanager.IRabbitMQManager) error
}
