package transaction

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	return c.JSON(http.StatusOK, "Transaction get transaction endpoint")
}
