package account

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AccountHandler) CreateMerchantAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, "Account create merchant API endpoint")
}
