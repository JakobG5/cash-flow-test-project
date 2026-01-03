package accountservice

import (
	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	ErrDuplicateEmail = errors.New("merchant with this email already exists")
)

type AccountService struct {
	queries *db.Queries
	logger  *loggermanager.Logger
	hashKey string
}

func NewAccountService(queries *db.Queries, logger *loggermanager.Logger, config *models.Config) IAccountService {
	return &AccountService{
		queries: queries,
		logger:  logger,
		hashKey: config.APIKeyHash,
	}
}

func (as *AccountService) CreateMerchant(name, email string) (*models.CreateMerchantResponse, error) {
	as.logger.Info("Starting merchant creation process", zap.String("email", email), zap.String("name", name))

	merchantID := as.generateMerchantID()
	apiKey := as.generateAPIKey()
	hashedAPIKey := as.hashAPIKey(apiKey)

	as.logger.Info("Generated merchant credentials", zap.String("merchant_id", merchantID))

	merchant, err := as.queries.CreateMerchant(context.Background(), &db.CreateMerchantParams{
		MerchantID: merchantID,
		Name:       name,
		Email:      email,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			as.logger.Warn("Duplicate email attempt", zap.String("email", email), zap.String("error", "unique constraint violation"))
			return nil, ErrDuplicateEmail
		}
		as.logger.Error("Failed to create merchant in database", zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	as.logger.Info("Merchant record created successfully", zap.String("merchant_id", merchantID))

	_, err = as.queries.CreateMerchantAPIKey(context.Background(), &db.CreateMerchantAPIKeyParams{
		MerchantID: merchant.ID,
		ApiKey:     hashedAPIKey,
		SecretKey:  as.generateSecretKey(),
	})
	if err != nil {
		as.logger.Error("Failed to create merchant API key", zap.String("merchant_id", merchant.ID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create merchant API key: %w", err)
	}

	as.logger.Info("Merchant API key created successfully", zap.String("merchant_id", merchant.ID.String()))

	as.logger.Info("Merchant creation completed successfully", zap.String("merchant_id", merchant.ID.String()), zap.String("email", email))

	return &models.CreateMerchantResponse{
		Status:     true,
		MerchantID: merchantID,
		Name:       name,
		Email:      email,
		APIKey:     apiKey,
		Message:    "Merchant created successfully",
	}, nil
}

func (as *AccountService) GetMerchantByID(merchantID string) (*models.GetMerchantResponse, error) {
	as.logger.Info("Getting merchant by ID", zap.String("merchant_id", merchantID))

	merchant, err := as.queries.GetMerchantWithAPIKey(context.Background(), merchantID)
	if err != nil {
		as.logger.Error("Failed to get merchant by ID", zap.String("merchant_id", merchantID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("merchant not found")
	}

	// Handle null status fields
	merchantStatus := "unknown"
	if merchant.Status.Valid {
		merchantStatus = string(merchant.Status.MerchantStatus)
	}

	apiKeyStatus := "unknown"
	if merchant.ApiKeyStatus.Valid {
		apiKeyStatus = string(merchant.ApiKeyStatus.ApiKeyStatus)
	}

	// Handle null timestamp fields
	createdAt := ""
	if merchant.CreatedAt.Valid {
		createdAt = merchant.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	apiKeyCreatedAt := ""
	if merchant.ApiKeyCreatedAt.Valid {
		apiKeyCreatedAt = merchant.ApiKeyCreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	as.logger.Info("Merchant retrieved successfully", zap.String("merchant_id", merchantID), zap.String("email", merchant.Email))

	return &models.GetMerchantResponse{
		Status:         true,
		MerchantID:     merchant.MerchantID,
		Name:           merchant.Name,
		Email:          merchant.Email,
		MerchantStatus: merchantStatus,
		APIKey:         merchant.ApiKey,
		APIKeyStatus:   apiKeyStatus,
		CreatedAt:      createdAt,
		APIKeyCreated:  apiKeyCreatedAt,
		Message:        "Merchant retrieved successfully",
	}, nil
}

func (as *AccountService) GetMerchantByAPIKey(merchantID string) error {
	return nil
}

func (as *AccountService) generateAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "api_" + string(b)
}

func (as *AccountService) generateMerchantID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "CASM-" + string(b)
}

func (as *AccountService) generateSecretKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 64)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "sk_" + string(b)
}

func (as *AccountService) hashAPIKey(apiKey string) string {
	block, err := aes.NewCipher([]byte(as.hashKey)[:32])
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (as *AccountService) decryptAPIKey(hashedKey string) (string, error) {
	block, err := aes.NewCipher([]byte(as.hashKey)[:32])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(hashedKey)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
