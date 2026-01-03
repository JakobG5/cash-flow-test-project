package accountservice

import (
	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"
	"context"
	"errors"
	"fmt"

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
	as.logger.Info("=== MERCHANT CREATION START ===")
	as.logger.Info("Starting merchant creation process", zap.String("email", email), zap.String("name", name))

	merchantID := generateMerchantID()
	apiKey := generateAPIKey()
	as.logger.Info("Generated plain API key", zap.String("plain_api_key", maskAPIKey(apiKey)))

	hashedAPIKey := hashAPIKey(apiKey, as.hashKey)
	hashPreview := hashedAPIKey
	if len(hashedAPIKey) > 50 {
		hashPreview = hashedAPIKey[:50] + "..."
	}
	as.logger.Info("Generated hashed API key", zap.String("hashed_preview", hashPreview))

	encryptedAPIKey := encryptAPIKey(apiKey, as.hashKey)
	as.logger.Info("Generated encrypted API key for storage", zap.String("encrypted_preview", encryptedAPIKey[:50]+"..."))

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
		SecretKey:  encryptedAPIKey,
	})
	if err != nil {
		as.logger.Error("Failed to create merchant API key", zap.String("merchant_id", merchant.ID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create merchant API key: %w", err)
	}

	as.logger.Info("Merchant API key created successfully", zap.String("merchant_id", merchant.ID.String()))
	storedHashPreview := hashedAPIKey
	if len(hashedAPIKey) > 50 {
		storedHashPreview = hashedAPIKey[:50] + "..."
	}
	storedEncryptedPreview := encryptedAPIKey
	if len(encryptedAPIKey) > 50 {
		storedEncryptedPreview = encryptedAPIKey[:50] + "..."
	}

	as.logger.Info("=== DATABASE STORAGE DETAILS ===")
	as.logger.Info("Stored in api_key column (for auth)", zap.String("stored_hash", storedHashPreview))
	as.logger.Info("Stored in secret_key column (for retrieval)", zap.String("stored_encrypted", storedEncryptedPreview))
	as.logger.Info("Plain API key for user", zap.String("user_api_key", maskAPIKey(apiKey)))
	as.logger.Info("=== END DATABASE STORAGE ===")

	as.logger.Info("Merchant creation completed successfully", zap.String("merchant_id", merchant.ID.String()), zap.String("email", email))
	as.logger.Info("=== MERCHANT CREATION END ===")
	as.logger.Info("Returning API key to user", zap.String("returned_api_key", maskAPIKey(apiKey)))

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

	as.logger.Info("=== MERCHANT RETRIEVAL DETAILS ===")
	as.logger.Info("Retrieved from database", zap.String("merchant_id", merchantID))
	as.logger.Info("Stored api_key (hash)", zap.String("db_api_key", maskAPIKey(merchant.ApiKey)))
	as.logger.Info("Stored secret_key (encrypted)", zap.String("db_secret_key", maskAPIKey(merchant.SecretKey)))
	as.logger.Info("=== END RETRIEVAL DETAILS ===")

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

	decryptedAPIKey, err := decryptAPIKey(merchant.SecretKey, as.hashKey)
	if err != nil {
		as.logger.Error("Failed to decrypt API key", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve API key")
	}

	as.logger.Info("Merchant retrieved successfully", zap.String("merchant_id", merchantID), zap.String("email", merchant.Email))

	return &models.GetMerchantResponse{
		Status:         true,
		MerchantID:     merchant.MerchantID,
		Name:           merchant.Name,
		Email:          merchant.Email,
		MerchantStatus: merchantStatus,
		APIKey:         decryptedAPIKey,
		APIKeyStatus:   apiKeyStatus,
		CreatedAt:      createdAt,
		APIKeyCreated:  apiKeyCreatedAt,
		Message:        "Merchant retrieved successfully",
	}, nil
}

func (as *AccountService) GetMerchantByAPIKey(apiKey string) (*models.GetMerchantResponse, error) {
	as.logger.Info("=== API KEY AUTHENTICATION START ===")
	as.logger.Info("Received API key for authentication", zap.String("received_key", maskAPIKey(apiKey)))
	as.logger.Info("API key length", zap.Int("length", len(apiKey)))

	hashedAPIKey := hashAPIKey(apiKey, as.hashKey)
	authHashPreview := hashedAPIKey
	if len(hashedAPIKey) > 50 {
		authHashPreview = hashedAPIKey[:50] + "..."
	}
	as.logger.Info("Generated hash for lookup", zap.String("hashed_key_preview", authHashPreview))
	as.logger.Info("Full hashed key length", zap.Int("hashed_length", len(hashedAPIKey)))

	as.logger.Info("Executing database query with hashed key")
	merchant, err := as.queries.GetMerchantByAPIKeyValue(context.Background(), hashedAPIKey)

	if err != nil {
		as.logger.Error("Database query failed", zap.Error(err))
		as.logger.Warn("Invalid API key lookup failed", zap.String("api_key", maskAPIKey(apiKey)))
		as.logger.Info("=== API KEY AUTHENTICATION FAILED ===")
		return nil, fmt.Errorf("invalid API key")
	}

	as.logger.Info("Database query successful - merchant found", zap.String("merchant_id", merchant.MerchantID))
	as.logger.Info("=== API KEY AUTHENTICATION SUCCESS ===")

	merchantStatus := "unknown"
	if merchant.Status.Valid {
		merchantStatus = string(merchant.Status.MerchantStatus)
	}

	apiKeyStatus := "unknown"
	if merchant.ApiKeyStatus.Valid {
		apiKeyStatus = string(merchant.ApiKeyStatus.ApiKeyStatus)
	}

	createdAt := ""
	if merchant.CreatedAt.Valid {
		createdAt = merchant.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	apiKeyCreatedAt := ""
	if merchant.ApiKeyCreatedAt.Valid {
		apiKeyCreatedAt = merchant.ApiKeyCreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	as.logger.Info("Merchant found by API key", zap.String("merchant_id", merchant.MerchantID), zap.String("email", merchant.Email))

	decryptedAPIKey, err := decryptAPIKey(merchant.SecretKey, as.hashKey)
	if err != nil {
		as.logger.Error("Failed to decrypt API key for authentication response", zap.String("merchant_id", merchant.MerchantID), zap.Error(err))
		return nil, fmt.Errorf("failed to decrypt API key")
	}

	return &models.GetMerchantResponse{
		Status:         true,
		MerchantID:     merchant.MerchantID,
		Name:           merchant.Name,
		Email:          merchant.Email,
		MerchantStatus: merchantStatus,
		APIKey:         decryptedAPIKey,
		APIKeyStatus:   apiKeyStatus,
		CreatedAt:      createdAt,
		APIKeyCreated:  apiKeyCreatedAt,
		Message:        "Merchant found",
	}, nil
}
