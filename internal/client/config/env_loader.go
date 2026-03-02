package config

import "os"

// ParseEnv заполняет конфиг из SERVER_URL, если переменная задана.
func ParseEnv(cfg *Config) {
	if v := os.Getenv("SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
}
