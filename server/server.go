package server

import (
	"cash-flow-financial/internal/managers/dbmanager"
	logger "cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	"cash-flow-financial/internal/models"
	accountservice "cash-flow-financial/internal/services/account-service"
	checkoutservice "cash-flow-financial/internal/services/checkout-service"
	transactionservice "cash-flow-financial/internal/services/transaction-service"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	IDBManager          dbmanager.IDBManager
	IRabbitMQManager    rabbitmqmanager.IRabbitMQManager
	ICHECKOUTSERVICE    checkoutservice.ICheckoutService
	IACCOUNTSERVICE     accountservice.IAccountService
	ITRANSACTIONSERVICE transactionservice.ITransactionService
	echo                *echo.Echo
	config              *models.Config
	logger              *logger.Logger
}

func NewServer(cfg *models.Config, log *logger.Logger, checkoutSvc checkoutservice.ICheckoutService, accountSvc accountservice.IAccountService, transactionSvc transactionservice.ITransactionService, dbMgr dbmanager.IDBManager, rabbitMgr rabbitmqmanager.IRabbitMQManager) *Server {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server := &Server{
		IDBManager:          dbMgr,
		IRabbitMQManager:    rabbitMgr,
		ICHECKOUTSERVICE:    checkoutSvc,
		IACCOUNTSERVICE:     accountSvc,
		ITRANSACTIONSERVICE: transactionSvc,
		echo:                e,
		config:              cfg,
		logger:              log,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf(":%s", s.config.Server.Port)
	s.logger.Info("Starting server", zap.String("port", s.config.Server.Port))

	go func() {
		if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	s.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("Server gracefully stopped")
	return nil
}
