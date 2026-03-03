// Package config загружает настройки сервера из файла, переменных окружения и флагов.
package config

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/MaxRadzey/gog/internal/server/crypto"
)

// Config хранит настройки сервера.
type Config struct {
	Address       string `env:"SERVER_ADDRESS"` // адрес HTTP (например :8080)
	LogLevel      string `env:"LOG_LEVEL"`      // уровень логов (info, debug и т.д.)
	DatabaseDSN   string `env:"DATABASE_DSN"`   // DSN PostgreSQL; пусто — БД не используется
	SigningKey    string `env:"SECRET_KEY"`     // секрет для подписи куки
	EncryptionKey string `env:"ENCRYPTION_KEY"` // ключ шифрования секретов (32 байта)
	EnableHTTPS   bool   `env:"ENABLE_HTTPS"`   // использовать TLS (ListenAndServeTLS)
	TLSCertFile   string `env:"TLS_CERT_FILE"`  // путь к сертификату (для HTTPS)
	TLSKeyFile    string `env:"TLS_KEY_FILE"`   // путь к приватному ключу (для HTTPS)
}

// New возвращает конфиг: дефолты + файл (CONFIG / -c) + переменные окружения + флаги.
func New() *Config {
	cfg := &Config{
		Address:       "localhost:8080",
		LogLevel:      "info",
		DatabaseDSN:   "postgres://gog:gog@localhost:5433/gog_test",
		SigningKey:    "dev-signing-key-change-in-production",
		EncryptionKey: "dev-encryption-key-32bytes-long!",
		EnableHTTPS:   false,
		TLSCertFile:   "",
		TLSKeyFile:    "",
	}
	ParseFile(cfg)
	_ = cleanenv.ReadEnv(cfg)
	ParseFlags(cfg)
	return cfg
}

// Validate проверяет конфиг перед стартом приложения. При ошибке приложение должно завершиться.
func Validate(cfg *Config) error {
	if cfg.DatabaseDSN == "" {
		return errors.New("config: DATABASE_DSN is required")
	}
	if err := crypto.ValidateKey([]byte(cfg.EncryptionKey)); err != nil {
		return fmt.Errorf("config: encryption key: %w", err)
	}
	if cfg.EnableHTTPS {
		if cfg.TLSCertFile == "" {
			return errors.New("config: TLS_CERT_FILE is required when ENABLE_HTTPS is true")
		}
		if cfg.TLSKeyFile == "" {
			return errors.New("config: TLS_KEY_FILE is required when ENABLE_HTTPS is true")
		}
	}
	return nil
}
