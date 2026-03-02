// Package config — настройки клиента: URL сервера из env (SERVER_URL) и флагов (-server, -s).
package config

// Config хранит настройки клиента.
type Config struct {
	ServerURL string // базовый URL API (например http://localhost:8080)
}

// New создаёт конфиг с дефолтом, затем подставляет переменные окружения и флаги.
func New() *Config {
	cfg := &Config{
		ServerURL: "http://localhost:8080",
	}
	ParseEnv(cfg)
	ParseFlags(cfg)
	return cfg
}
