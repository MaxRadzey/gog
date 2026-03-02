// загрузка конфигурации из JSON (CONFIG, -c, -config).
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// fileConfig описывает поддерживаемые поля JSON-конфигурации.
type fileConfig struct {
	Address     *string `json:"server_address"`
	DatabaseDSN *string `json:"database_dsn"`
	LogLevel    *string `json:"log_level"`
	SigningKey  *string `json:"secret_key"`
}

// ParseFile заполняет config значениями из JSON-файла.
func ParseFile(config *Config) {
	path := getConfigPath()
	if path == "" {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open config file %q: %v\n", path, err)
		return
	}
	defer f.Close()

	var fc fileConfig
	if err := json.NewDecoder(f).Decode(&fc); err != nil {
		fmt.Fprintf(os.Stderr, "cannot decode config file %q: %v\n", path, err)
		return
	}

	if fc.Address != nil {
		config.Address = *fc.Address
	}
	if fc.DatabaseDSN != nil {
		config.DatabaseDSN = *fc.DatabaseDSN
	}
	if fc.LogLevel != nil {
		config.LogLevel = *fc.LogLevel
	}
	if fc.SigningKey != nil {
		config.SigningKey = *fc.SigningKey
	}

}

// getConfigPath возвращает путь к JSON-конфигу из CONFIG или флагов -c/-config.
func getConfigPath() string {
	if envPath := os.Getenv("CONFIG"); envPath != "" {
		return envPath
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-c" || arg == "-config" {
			if i+1 < len(args) {
				return args[i+1]
			}
			continue
		}

		if strings.HasPrefix(arg, "-c=") {
			return strings.TrimPrefix(arg, "-c=")
		}
		if strings.HasPrefix(arg, "-config=") {
			return strings.TrimPrefix(arg, "-config=")
		}
	}

	return ""
}
