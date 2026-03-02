package config

import "flag"

// ParseFlags читает флаги -a (адрес), -l (уровень логов), -d (DSN); приоритет выше env.
func ParseFlags(config *Config) {
	flag.StringVar(&config.Address, "a", config.Address, "address and port to run HTTP server")
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "log level")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "database connection string")

	flag.Parse()
}
