// Package config загружает настройки клиента из флагов и переменных окружения.
package config

// Config — настройки клиента (URL сервера).
type Config struct {
	ServerURL string // базовый URL сервера (например http://localhost:8080)
}

// New возвращает конфиг с дефолтами, затем подставляет env и флаги.
func New() *Config {
	cfg := &Config{
		ServerURL: "http://localhost:8080",
	}
	ParseEnv(cfg)
	ParseFlags(cfg)
	return cfg
}
