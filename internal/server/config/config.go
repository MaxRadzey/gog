// Package config загружает настройки сервера из файла, переменных окружения и флагов.
package config

import "github.com/ilyakaznacheev/cleanenv"

// Config хранит настройки сервера.
type Config struct {
	Address       string `env:"SERVER_ADDRESS"` // адрес HTTP (например :8080)
	LogLevel      string `env:"LOG_LEVEL"`      // уровень логов (info, debug и т.д.)
	DatabaseDSN   string `env:"DATABASE_DSN"`   // DSN PostgreSQL; пусто — БД не используется
	SigningKey    string `env:"SECRET_KEY"`     // секрет для подписи куки
	EncryptionKey string `env:"ENCRYPTION_KEY"` // ключ шифрования секретов (32 байта)
}

// New возвращает конфиг: дефолты + файл (CONFIG / -c) + переменные окружения + флаги.
func New() *Config {
	cfg := &Config{
		Address:       "localhost:8080",
		LogLevel:      "info",
		DatabaseDSN:   "postgres://shortener:shortener@localhost:5432/shortener",
		SigningKey:    "dev-signing-key-change-in-production",
		EncryptionKey: "dev-encryption-key-32bytes-long!",
	}
	ParseFile(cfg)
	_ = cleanenv.ReadEnv(cfg)
	ParseFlags(cfg)
	return cfg
}
