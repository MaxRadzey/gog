// Package logger — общий логгер на zap для клиента и сервера.
package logger

import (
	"go.uber.org/zap"
)

// Log — глобальный экземпляр логгера. До вызова Initialize — no-op.
var Log *zap.Logger = zap.NewNop()

// Initialize поднимает логгер с заданным уровнем (debug, info, warn, error).
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
