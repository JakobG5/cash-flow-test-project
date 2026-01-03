package account

import (
	"errors"
	"net/http"

	"cash-flow-financial/internal/models"
	accountservice "cash-flow-financial/internal/services/account-service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *AccountHandler) CreateMerchantAPI(c echo.Context) error {
	h.logger.Info("CreateMerchantAPI called")

	var req models.CreateMerchantRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("CreateMerchantAPI failed: invalid request format")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Status: false, Error: "invalid request format"})
	}

	if validationErrors := h.validateCreateMerchantRequest(req); len(validationErrors) > 0 {
		h.logger.Warn("CreateMerchantAPI failed: validation errors", zap.Strings("errors", validationErrors))
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Error:   "validation failed",
			Details: validationErrors,
		})
	}

	response, err := h.accountService.CreateMerchant(req.Name, req.Email)
	if err != nil {
		if errors.Is(err, accountservice.ErrDuplicateEmail) {
			h.logger.Warn("CreateMerchantAPI failed: duplicate email", zap.String("email", req.Email))
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status:  false,
				Error:   "duplicate email",
				Details: []string{"A merchant with this email already exists"},
			})
		}
		h.logger.Error("CreateMerchantAPI failed: internal error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status: false,
			Error:  "internal server error",
		})
	}

	h.logger.Info("CreateMerchantAPI successful", zap.String("merchant_id", response.MerchantID))
	return c.JSON(http.StatusCreated, response)
}
