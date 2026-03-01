package crypto

import "errors"

// ErrKeyLength — ключ должен быть 32 байта (AES-256).
var ErrKeyLength = errors.New("encryption key must be 32 bytes")

// ErrDecrypt — не удалось расшифровать (неверный ключ или повреждённые данные).
var ErrDecrypt = errors.New("decryption failed")
