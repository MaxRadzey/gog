// Package app собирает логгер, хранилище, сервисы и запускает HTTP-сервер с graceful shutdown.
package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaxRadzey/gog/internal/logger"
	"github.com/MaxRadzey/gog/internal/server/api"
	"github.com/MaxRadzey/gog/internal/server/api/handler"
	"github.com/MaxRadzey/gog/internal/server/config"
	"github.com/MaxRadzey/gog/internal/server/service"
	"github.com/MaxRadzey/gog/internal/server/storage"
	"go.uber.org/zap"
)

const shutdownTimeout = 30 * time.Second

// App — приложение: конфиг, HTTP-сервер и хранилище.
type App struct {
	config  *config.Config
	server  *http.Server
	storage *storage.Storage
}

// New создаёт приложение: инициализирует логгер, хранилище, сервисы и HTTP-сервер.
func New(cfg *config.Config) (*App, error) {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}

	storageResult, err := storage.InitializeStorage(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	services := service.NewServices(storageResult, cfg.EncryptionKey)
	h := handler.New(services, cfg.SigningKey)
	router := api.SetupRouter(h)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	return &App{
		config:  cfg,
		server:  server,
		storage: storageResult,
	}, nil
}

// Run запускает сервер, ждёт сигнал завершения (SIGTERM, SIGINT, SIGQUIT) и корректно останавливает приложение.
func (a *App) Run() error {
	a.startServer()
	<-a.shutdownSignal()
	return a.gracefulShutdown()
}

func (a *App) startServer() {
	logger.Log.Info("Starting HTTP server", zap.String("address", a.config.Address))
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("HTTP server error", zap.Error(err))
		}
	}()
}

func (a *App) shutdownSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return ch
}

func (a *App) gracefulShutdown() error {
	logger.Log.Info("Shutdown signal received, draining connections...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		logger.Log.Error("server shutdown error", zap.Error(err))
		return err
	}

	a.closeStorage()
	logger.Log.Info("Server stopped gracefully")
	return nil
}

func (a *App) closeStorage() {
	if a.storage != nil {
		if err := a.storage.Close(); err != nil {
			logger.Log.Error("storage close error", zap.Error(err))
		} else {
			logger.Log.Info("Storage closed")
		}
	}
}
