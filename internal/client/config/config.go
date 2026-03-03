// Package config загружает настройки клиента из переменных окружения и флагов.
package config

import "github.com/ilyakaznacheev/cleanenv"

// Config хранит настройки клиента.
type Config struct {
	ServerURL string `env:"SERVER_URL"` // базовый URL API (например http://localhost:8080)
}

// New возвращает конфиг: дефолт + переменные окружения + флаги.
func New() *Config {
	cfg := &Config{
		ServerURL: "http://localhost:8080",
	}
	_ = cleanenv.ReadEnv(cfg)
	ParseFlags(cfg)
	return cfg
}
