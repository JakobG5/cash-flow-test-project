package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/configmanager"
	"cash-flow-financial/internal/managers/dbmanager"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	accountservice "cash-flow-financial/internal/services/account-service"
	"cash-flow-financial/internal/services/callback"
	checkoutservice "cash-flow-financial/internal/services/checkout-service"
	transactionservice "cash-flow-financial/internal/services/transaction-service"
	"cash-flow-financial/server"
	"cash-flow-financial/worker"

	_ "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := configmanager.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	logger := loggermanager.NewLogger(cfg.Logger.Level)

	dbManager, err := dbmanager.NewDBManager(&cfg.Database)
	if err != nil {
		panic("Failed to initialize database manager: " + err.Error())
	}
	defer dbManager.Close()

	rabbitManager, err := rabbitmqmanager.NewRabbitMQManager(&cfg.RabbitMQ)
	if err != nil {
		panic("Failed to initialize RabbitMQ manager: " + err.Error())
	}
	defer rabbitManager.Close()

	queries := db.New(dbManager.GetDB())

	checkoutService := checkoutservice.NewCheckoutService(queries, logger, rabbitManager)
	accountService := accountservice.NewAccountService(queries, logger, cfg)
	transactionService := transactionservice.NewTransactionService(queries, logger)

	// Initialize callback service
	callbackService := callback.NewCallbackService(logger, cfg)

	// Initialize worker
	paymentWorker := worker.NewWorker(queries, rabbitManager, logger, callbackService)

	srv := server.NewServer(cfg, logger, checkoutService, accountService, transactionService, dbManager, rabbitManager)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the payment worker
	if err := paymentWorker.Start(ctx); err != nil {
		logger.Fatal("Worker failed to start", zap.Error(err))
	}

	if err := srv.Start(ctx); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
