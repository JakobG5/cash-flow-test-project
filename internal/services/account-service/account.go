package accountservice

type IAccountService interface {
	CreateMerchant(name, email string) error
	GetMerchantByAPIKey(merchantID string) error
}
