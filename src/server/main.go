package main

// @title Song API
// @version 1.0
// @description API для управления песнями
// @host localhost:8085
// @BasePath /

import (
	"EffectiveMobileTestTask/internal/config"
	"EffectiveMobileTestTask/internal/externalApiClient"
	"EffectiveMobileTestTask/internal/handlers"
	"EffectiveMobileTestTask/internal/libaryDB"
	"EffectiveMobileTestTask/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Init logger
	logger := NewLogger()

	logger.Info("Logger initialized")

	cfg := config.NewConfig()
	logger.Info("Configuration initialized")

	logger.Info("Waiting for postgres container...")
	time.Sleep(time.Second * 5)

	db := libaryDB.NewDB(logger, cfg)
	logger.Info("DB initialized")

	apiClient := externalApiClient.NewAPIClient(logger)
	logger.Info("External API client initialized")

	handler := handlers.NewHandler(db, apiClient, logger)
	logger.Info("Handler initialized")

	srv := server.NewServer(logger, handler)
	logger.Info("Server initialized")

	go srv.Start(cfg.ServerConfig.Port)
	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	srv.Stop(time.Second * 5)
	db.Shutdown()
	logger.Info("Shut down server")
}

func NewLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}
