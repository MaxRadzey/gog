package config

import "os"

// ParseEnv подставляет в config значения из переменных окружения.
func ParseEnv(cfg *Config) {
	if v := os.Getenv("SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
}
