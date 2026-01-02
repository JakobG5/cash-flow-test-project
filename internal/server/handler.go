package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "cash-flow-financial",
	})
}

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
		"service": "cash-flow-financial",
	})
}
