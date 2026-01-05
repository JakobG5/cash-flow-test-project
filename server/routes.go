package server

import (
	"cash-flow-financial/server/handlers/account"
	"cash-flow-financial/server/handlers/checkout"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) setupRoutes() {
	// Health check
	s.echo.GET("/health", s.healthCheck)

	// Swagger documentation
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	checkoutHandler := checkout.NewCheckoutHandler(s.ICHECKOUTSERVICE, s.IACCOUNTSERVICE, s.config, s.logger, s.IRabbitMQManager)
	accountHandler := account.NewAccountHandler(s.IACCOUNTSERVICE, s.config, s.logger)

	apiV1 := s.echo.Group("/cashflow_test/v1")

	// Checkout routes
	apiV1.POST("/checkout/create-intent", checkoutHandler.CreateIntent)

	// Account routes
	apiV1.POST("/account/create-merchant", accountHandler.CreateMerchantAPI)
	apiV1.GET("/account/merchant", accountHandler.GetMerchantAPI) // Requires merchant_id query param, returns merchant details, balances, and transactions
}

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"status":  "healthy",
		"service": "cash-flow-financial",
	})
}
