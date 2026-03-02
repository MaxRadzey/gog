package config

import "flag"

// ParseFlags парсит флаги; приоритет над env.
func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.ServerURL, "server", cfg.ServerURL, "server base URL")
	flag.StringVar(&cfg.ServerURL, "s", cfg.ServerURL, "server base URL (short)")
	flag.Parse()
}
