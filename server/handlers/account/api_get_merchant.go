package account

import (
	"net/http"

	"cash-flow-financial/internal/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// GetMerchantAPI retrieves merchant details, balances, and transaction history
// @Summary Get Merchant Details
// @Description Retrieves merchant information including balances across currencies and recent transactions
// @Tags Merchant
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID (e.g., CASM-ABC123)"
// @Success 200 {object} models.GetMerchantResponse "Merchant details retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Missing merchant_id parameter"
// @Failure 404 {object} models.ErrorResponse "Merchant not found"
// @Router /account/merchant [get]
func (h *AccountHandler) GetMerchantAPI(c echo.Context) error {
	merchantID := c.QueryParam("merchant_id")
	h.logger.Info("GetMerchantAPI called", zap.String("merchant_id", merchantID))

	if merchantID == "" {
		h.logger.Warn("GetMerchantAPI failed: missing merchant_id parameter")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: false,
			Error:  "merchant_id query parameter is required",
		})
	}

	response, err := h.accountService.GetMerchantByID(merchantID)
	if err != nil {
		h.logger.Warn("GetMerchantAPI failed: merchant not found", zap.String("merchant_id", merchantID))
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: false,
			Error:  "merchant not found",
		})
	}

	h.logger.Info("GetMerchantAPI successful", zap.String("merchant_id", merchantID))
	return c.JSON(http.StatusOK, response)
}
