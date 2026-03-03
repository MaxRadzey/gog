package crypto

import "errors"

// ErrKeyLength возвращается, если ключ не 32 байта.
var ErrKeyLength = errors.New("encryption key must be 32 bytes")

// ErrDecrypt возвращается при ошибке расшифровки (неверный ключ или битые данные).
var ErrDecrypt = errors.New("decryption failed")
