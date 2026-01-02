package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cash-flow-financial/internal/managers/configmanager"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/server"

	"go.uber.org/zap"
)

func main() {
	cfg, err := configmanager.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	log := loggermanager.NewLogger(cfg.Logger.Level)

	srv := server.NewServer(cfg, log)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}
