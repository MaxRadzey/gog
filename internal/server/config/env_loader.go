package config

import (
	"os"
)

// ParseEnv заполняет конфиг из SERVER_ADDRESS, LOG_LEVEL, DATABASE_DSN, SECRET_KEY, ENCRYPTION_KEY (если заданы).
func ParseEnv(config *Config) {
	if Address := os.Getenv("SERVER_ADDRESS"); Address != "" {
		config.Address = Address
	}
	if LogLevel := os.Getenv("LOG_LEVEL"); LogLevel != "" {
		config.LogLevel = LogLevel
	}
	if DatabaseDSN := os.Getenv("DATABASE_DSN"); DatabaseDSN != "" {
		config.DatabaseDSN = DatabaseDSN
	}
	if v := os.Getenv("SECRET_KEY"); v != "" {
		config.SigningKey = v
	}
	if v := os.Getenv("ENCRYPTION_KEY"); v != "" {
		config.EncryptionKey = v
	}
}
