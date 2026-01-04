package account

import (
	"net/http"

	"cash-flow-financial/internal/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

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
