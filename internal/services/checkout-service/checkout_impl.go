package checkoutservice

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	"cash-flow-financial/internal/models"

	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
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

func (cs *CheckoutService) CreatePaymentIntent(merchantID string, req models.CreatePaymentIntentRequest) (*models.CreatePaymentIntentResponse, error) {
	cs.logger.Info("Creating payment intent", zap.String("merchant_id", merchantID), zap.Float64("amount", req.Amount), zap.String("nonce", req.Nonce))

	// Get merchant UUID from merchant_id string
	merchant, err := cs.queries.GetMerchantByMerchantID(context.Background(), merchantID)
	if err != nil {
		cs.logger.Error("Failed to get merchant", zap.String("merchant_id", merchantID), zap.Error(err))
		return nil, errors.New("invalid merchant")
	}

	// Check for idempotency - see if payment intent with same nonce already exists for this merchant
	existingIntent, err := cs.queries.GetPaymentIntentByNonce(context.Background(), &db.GetPaymentIntentByNonceParams{
		MerchantID: merchant.ID,
		Nonce:      req.Nonce,
	})
	if err == nil {
		// Payment intent already exists, return it (idempotent behavior)
		cs.logger.Info("Payment intent already exists for nonce, returning existing", zap.String("existing_payment_intent_id", existingIntent.PaymentIntentID), zap.String("nonce", req.Nonce))

		amount, _ := strconv.ParseFloat(existingIntent.Amount, 64)
		var description string
		if existingIntent.Description.Valid {
			description = existingIntent.Description.String
		}
		paymentStatus := string(existingIntent.Status.PaymentStatus)
		if !existingIntent.Status.Valid {
			paymentStatus = "unknown"
		}

		return &models.CreatePaymentIntentResponse{
			Status:          true,
			PaymentIntentID: existingIntent.PaymentIntentID,
			Amount:          amount,
			Currency:        existingIntent.Currency,
			PaymentStatus:   paymentStatus,
			Description:     description,
			CreatedAt:       existingIntent.CreatedAt.Time,
			ExpiresAt:       existingIntent.ExpiresAt.Time,
			Message:         "Payment intent already exists",
		}, nil
	}

	// Generate new payment intent ID
	paymentIntentID := cs.generatePaymentIntentID()

	// Prepare metadata as JSON
	var metadata pqtype.NullRawMessage
	if req.Metadata != nil {
		jsonData, err := json.Marshal(req.Metadata)
		if err != nil {
			cs.logger.Error("Failed to marshal metadata", zap.Error(err))
			return nil, errors.New("invalid metadata")
		}
		metadata = pqtype.NullRawMessage{
			RawMessage: jsonData,
			Valid:      true,
		}
	}

	// Create new payment intent
	intent, err := cs.queries.CreatePaymentIntent(context.Background(), &db.CreatePaymentIntentParams{
		PaymentIntentID: paymentIntentID,
		MerchantID:      merchant.ID,
		Amount:          strconv.FormatFloat(req.Amount, 'f', 2, 64),
		Currency:        req.Currency,
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
		CallbackUrl:     req.CallbackURL,
		Nonce:           req.Nonce,
		Metadata:        metadata,
	})
	if err != nil {
		cs.logger.Error("Failed to create payment intent in database", zap.Error(err))
		return nil, errors.New("failed to create payment intent")
	}

	cs.logger.Info("Payment intent created successfully", zap.String("payment_intent_id", paymentIntentID), zap.String("nonce", req.Nonce))

	// TODO: Publish to RabbitMQ for async processing
	// cs.rabbitManager.PublishPaymentMessage(paymentIntentID, ...)

	amount, _ := strconv.ParseFloat(intent.Amount, 64)
	var description string
	if intent.Description.Valid {
		description = intent.Description.String
	}
	paymentStatus := string(intent.Status.PaymentStatus)
	if !intent.Status.Valid {
		paymentStatus = "unknown"
	}

	return &models.CreatePaymentIntentResponse{
		Status:          true,
		PaymentIntentID: intent.PaymentIntentID,
		Amount:          amount,
		Currency:        intent.Currency,
		PaymentStatus:   paymentStatus,
		Description:     description,
		CreatedAt:       intent.CreatedAt.Time,
		ExpiresAt:       intent.ExpiresAt.Time,
		Message:         "Payment intent created successfully",
	}, nil
}

func (cs *CheckoutService) generatePaymentIntentID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "PI-" + string(b)
}
