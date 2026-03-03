// Package crypto — шифрование и расшифровка данных (AES-256-GCM, ключ 32 байта).
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const aesKeyLen = 32

// ValidateKey проверяет, что ключ шифрования имеет длину 32 байта. Вызывать при старте приложения.
func ValidateKey(key []byte) error {
	if len(key) != aesKeyLen {
		return ErrKeyLength
	}
	return nil
}

// Encrypt шифрует plaintext ключом key (32 байта). Возвращает ciphertext с префиксом nonce.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != aesKeyLen {
		return nil, ErrKeyLength
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt расшифровывает ciphertext (с префиксом nonce) ключом key. Возвращает plaintext или ErrDecrypt.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != aesKeyLen {
		return nil, ErrKeyLength
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrDecrypt
	}
	nonce, ciphertextBody := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertextBody, nil)
	if err != nil {
		return nil, ErrDecrypt
	}
	return plain, nil
}
