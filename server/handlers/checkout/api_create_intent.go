package checkout

import (
	"net/http"
	"strings"

	"cash-flow-financial/internal/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CreateIntent creates a new payment intent for processing
// @Summary Create Payment Intent
// @Description Creates a payment intent that will be processed asynchronously. The payment gateway charges 1% fee and only accepts ETB and USD currencies.
// @Tags Payment
// @Accept json
// @Produce json
// @Param X-API-KEY header string true "Merchant API Key"
// @Param request body models.CreatePaymentIntentRequest true "Payment intent creation request"
// @Success 201 {object} models.CreatePaymentIntentResponse "Payment intent created successfully"
// @Failure 400 {object} models.ErrorResponse "Validation error"
// @Failure 401 {object} models.ErrorResponse "Invalid or missing API key"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /checkout/create-intent [post]
func (h *CheckoutHandler) CreateIntent(c echo.Context) error {
	h.logger.Info("CreateIntent called")

	// Extract API key from header
	apiKey := c.Request().Header.Get("X-API-KEY")
	if apiKey == "" {
		h.logger.Warn("CreateIntent failed: missing X-API-KEY header")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status: false,
			Error:  "X-API-KEY header is required",
		})
	}

	// Trim any whitespace
	apiKey = strings.TrimSpace(apiKey)

	// Validate API key and get merchant
	merchantResponse, err := h.accountService.GetMerchantByAPIKey(apiKey)
	if err != nil {
		if strings.Contains(err.Error(), "invalid API key") {
			h.logger.Warn("CreateIntent failed: invalid API key", zap.String("api_key", maskAPIKey(apiKey)))
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Status: false,
				Error:  "invalid API key",
			})
		}
		h.logger.Error("CreateIntent failed: merchant lookup error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status: false,
			Error:  "internal server error",
		})
	}

	// Parse request payload
	var req models.CreatePaymentIntentRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("CreateIntent failed: invalid request format")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: false,
			Error:  "invalid request format",
		})
	}

	// Validate request
	if validationErrors := h.validateCreatePaymentIntentRequest(req); len(validationErrors) > 0 {
		h.logger.Warn("CreateIntent failed: validation errors", zap.Strings("errors", validationErrors))
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Error:   "validation failed",
			Details: validationErrors,
		})
	}

	// Create payment intent
	response, err := h.checkoutService.CreatePaymentIntent(merchantResponse.MerchantID, req)
	if err != nil {
		h.logger.Error("CreateIntent failed: payment intent creation error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status: false,
			Error:  "failed to create payment intent",
		})
	}

	h.logger.Info("CreateIntent successful", zap.String("payment_intent_id", response.PaymentIntentID), zap.String("merchant_id", merchantResponse.MerchantID))
	return c.JSON(http.StatusCreated, response)
}
