package accountservice

import (
	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
)

type AccountService struct {
	queries *db.Queries
	logger  *loggermanager.Logger
}

func NewAccountService(queries *db.Queries, logger *loggermanager.Logger) IAccountService {
	return &AccountService{
		queries: queries,
		logger:  logger,
	}
}

func (as *AccountService) CreateMerchant(name, email string) error {
	return nil
}

func (as *AccountService) GetMerchantByAPIKey(merchantID string) error {
	return nil
}
