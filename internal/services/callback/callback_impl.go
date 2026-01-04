package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"

	"go.uber.org/zap"
)

type CallbackService struct {
	logger *loggermanager.Logger
	config *models.Config
	client *http.Client
}

func NewCallbackService(logger *loggermanager.Logger, config *models.Config) ICallbackService {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &CallbackService{
		logger: logger,
		config: config,
		client: client,
	}
}

func (cs *CallbackService) SendCallback(callbackURL string, request CallbackRequest) error {
	cs.logger.Info("Sending callback to merchant",
		zap.String("callback_url", callbackURL),
		zap.String("payment_intent_id", request.PaymentIntentID),
		zap.String("merchant_id", request.MerchantID),
		zap.String("status", request.Status))

	requestBody, err := json.Marshal(request)
	if err != nil {
		cs.logger.Error("Failed to marshal callback request", zap.Error(err))
		return fmt.Errorf("failed to marshal callback request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(requestBody))
	if err != nil {
		cs.logger.Error("Failed to create callback HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create callback HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "CashFlow-Financial/1.0")

	resp, err := cs.client.Do(httpReq)
	if err != nil {
		cs.logger.Error("Failed to send callback HTTP request",
			zap.String("callback_url", callbackURL),
			zap.Error(err))
		return fmt.Errorf("failed to send callback HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		cs.logger.Warn("Callback request failed with non-2xx status",
			zap.String("callback_url", callbackURL),
			zap.Int("status_code", resp.StatusCode),
			zap.String("payment_intent_id", request.PaymentIntentID))
		return nil
	}

	cs.logger.Info("Callback sent successfully",
		zap.String("callback_url", callbackURL),
		zap.Int("status_code", resp.StatusCode),
		zap.String("payment_intent_id", request.PaymentIntentID))

	return nil
}
