package command

import (
	"encoding/json"
	"errors"
	"strings"
)

// parseKeyValueArgs превращает аргументы вида "key=value" в JSON.
// Набор полей произвольный и зависит от типа секрета (login, card, text, file и т.д.).
func parseKeyValueArgs(args []string) (json.RawMessage, error) {
	if len(args) == 0 {
		return []byte("{}"), nil
	}
	m := make(map[string]string, len(args))
	for _, a := range args {
		idx := strings.Index(a, "=")
		if idx == -1 {
			return nil, errors.New("expected key=value, got " + a)
		}
		key := strings.TrimSpace(a[:idx])
		value := strings.TrimSpace(a[idx+1:])
		if key == "" {
			return nil, errors.New("empty key in key=value")
		}
		m[key] = value
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
