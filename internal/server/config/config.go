// Package config — настройки сервера: адрес, БД, уровень логов, ключи; загрузка из env, флагов и JSON-файла.
package config

// Config хранит настройки сервера.
type Config struct {
	Address       string // адрес HTTP (например :8080)
	LogLevel      string // уровень логов (info, debug и т.д.)
	DatabaseDSN   string // DSN PostgreSQL; пусто — БД не используется
	SigningKey    string // секрет для подписи куки
	EncryptionKey string // ключ шифрования секретов (32 байта)
}

// New возвращает конфиг с дефолтами.
func New() *Config {
	return &Config{
		Address:       "localhost:8080",
		LogLevel:      "info",
		DatabaseDSN:   "postgres://shortener:shortener@localhost:5432/shortener",
		SigningKey:    "dev-signing-key-change-in-production",
		EncryptionKey: "dev-encryption-key-32bytes-long!",
	}
}
