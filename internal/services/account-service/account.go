package accountservice

import "cash-flow-financial/internal/models"

type IAccountService interface {
	CreateMerchant(name, email string) (*models.CreateMerchantResponse, error)
	GetMerchantByID(merchantID string) (*models.GetMerchantResponse, error)
	GetMerchantByAPIKey(merchantID string) error
}
