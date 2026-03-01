// Package config загружает настройки приложения из флагов и переменных окружения.
package config

// Config — настройки приложения (адрес, БД, файл, логи, аудит, ключ подписи).
type Config struct {
	Address     string // адрес HTTP-сервера (например :8080)
	LogLevel    string // уровень логов (info, debug и т.д.)
	DatabaseDSN string // строка подключения к PostgreSQL; пусто — БД не используется
	SigningKey  string // SigningKey — секрет для подписи куки (в проде задавать через SECRET_KEY)

}

// New возвращает конфиг с дефолтными значениями.
func New() *Config {
	return &Config{
		Address:     "localhost:8080",
		LogLevel:    "info",
		DatabaseDSN: "postgres://shortener:shortener@localhost:5432/shortener",
		SigningKey:  "dev-signing-key-change-in-production",
	}
}
