package config

import "flag"

// ParseFlags читает флаги -server и -s; приоритет выше, чем у переменных окружения.
func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.ServerURL, "server", cfg.ServerURL, "server base URL")
	flag.StringVar(&cfg.ServerURL, "s", cfg.ServerURL, "server base URL (short)")
	flag.Parse()
}
