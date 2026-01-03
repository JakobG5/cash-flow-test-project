package server

import (
	"cash-flow-financial/server/handlers/account"
	"cash-flow-financial/server/handlers/checkout"
	"cash-flow-financial/server/handlers/transaction"
)

func (s *Server) setupRoutes() {
	checkoutHandler := checkout.NewCheckoutHandler(s.ICHECKOUTSERVICE, s.config, s.logger, s.IRabbitMQManager)
	accountHandler := account.NewAccountHandler(s.IACCOUNTSERVICE, s.config, s.logger)
	transactionHandler := transaction.NewTransactionHandler(s.ITRANSACTIONSERVICE, s.config, s.logger)
	s.echo.POST("/checkout/create-intent", checkoutHandler.CreateIntent)
	s.echo.POST("/account/create-merchant", accountHandler.CreateMerchantAPI)
	s.echo.GET("/transaction/get-payment-status", transactionHandler.GetTransaction)
}
