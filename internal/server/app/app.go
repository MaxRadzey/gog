// Package app — точка входа сервера: инициализация логгера, БД, сервисов и HTTP-сервера с graceful shutdown.
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

// App хранит конфиг, HTTP-сервер и хранилище БД.
type App struct {
	config  *config.Config
	server  *http.Server
	storage *storage.Storage
}

// New поднимает логгер, подключает БД, создаёт сервисы и роутер, возвращает приложение.
func New(cfg *config.Config) (*App, error) {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}

	if err := config.Validate(cfg); err != nil {
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
		Addr:              cfg.Address,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{
		config:  cfg,
		server:  server,
		storage: storageResult,
	}, nil
}

// Run запускает HTTP-сервер, ждёт SIGTERM/SIGINT/SIGQUIT и выполняет graceful shutdown.
func (a *App) Run() error {
	a.startServer()
	<-a.shutdownSignal()
	return a.gracefulShutdown()
}

func (a *App) startServer() {
	logger.Log.Info("Starting HTTP server", zap.String("address", a.config.Address))
	go func() {
		var err error
		if a.config.EnableHTTPS {
			logger.Log.Info("HTTPS mode enabled")
			err = a.server.ListenAndServeTLS(a.config.TLSCertFile, a.config.TLSKeyFile)
		} else {
			err = a.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
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
