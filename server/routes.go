package server

import (
	"cash-flow-financial/server/handlers/account"
	"cash-flow-financial/server/handlers/checkout"
	"cash-flow-financial/server/handlers/transaction"

	"github.com/labstack/echo/v4"
)

func (s *Server) setupRoutes() {
	s.echo.GET("/health", s.healthCheck)

	checkoutHandler := checkout.NewCheckoutHandler(s.ICHECKOUTSERVICE, s.IACCOUNTSERVICE, s.config, s.logger, s.IRabbitMQManager)
	accountHandler := account.NewAccountHandler(s.IACCOUNTSERVICE, s.config, s.logger)
	transactionHandler := transaction.NewTransactionHandler(s.ITRANSACTIONSERVICE, s.config, s.logger)

	apiV1 := s.echo.Group("/cashflow_test/v1")

	// Checkout routes
	apiV1.POST("/checkout/create-intent", checkoutHandler.CreateIntent)

	// Account routes
	apiV1.POST("/account/create-merchant", accountHandler.CreateMerchantAPI)
	apiV1.GET("/account/merchant", accountHandler.GetMerchantAPI)

	// Transaction routes
	apiV1.GET("/transaction/get-payment-status", transactionHandler.GetTransaction)
}

// healthCheck provides a simple health check endpoint
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"status":  "healthy",
		"service": "cash-flow-financial",
	})
}
