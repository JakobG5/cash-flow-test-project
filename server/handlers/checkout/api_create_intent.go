package checkout

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *CheckoutHandler) CreateIntent(c echo.Context) error {
	return c.JSON(http.StatusOK, "Checkout intnet endpoint")
}
