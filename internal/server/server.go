package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cash-flow-financial/internal/config"
	"cash-flow-financial/internal/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	logger *logger.Logger
}

func NewServer(cfg *config.Config, log *logger.Logger) *Server {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	SetupRoutes(e)

	return &Server{
		echo:   e,
		config: cfg,
		logger: log,
	}
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
