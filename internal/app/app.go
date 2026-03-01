package app

import (
	"github.com/MaxRadzey/gog/internal/config"
	"github.com/MaxRadzey/gog/internal/logger"
	"github.com/MaxRadzey/gog/internal/server_http"
	"github.com/MaxRadzey/gog/internal/server_http/handler"
	"github.com/MaxRadzey/gog/internal/service"
	"github.com/MaxRadzey/gog/internal/storage"
	"go.uber.org/zap"
)

// Run запускает http сервер.
func Run(cfg *config.Config) error {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return err
	}

	storageResult, err := storage.InitializeStorage(cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer storageResult.Close()

	services := service.NewServices(storageResult)

	h := handler.New(services, cfg.SigningKey)
	r := server_http.SetupRouter(h)
	logger.Log.Info("Starting HTTP server", zap.String("address", cfg.Address))
	return r.Run(cfg.Address)
}
