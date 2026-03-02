// Package logger предоставляет общий логгер (zap) для клиента и сервера.
package logger

import (
	"go.uber.org/zap"
)

// Log — глобальный логгер. По умолчанию no-op; инициализируется через Initialize.
var Log *zap.Logger = zap.NewNop()

// Initialize настраивает логгер по уровню (debug, info, warn, error).
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
