package checkoutservice

import "cash-flow-financial/internal/models"

type ICheckoutService interface {
	CreatePaymentIntent(merchantID string, req models.CreatePaymentIntentRequest) (*models.CreatePaymentIntentResponse, error)
}
