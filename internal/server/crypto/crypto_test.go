package crypto

import (
	"bytes"
	"errors"
	"testing"
)

func mustKey() []byte {
	return []byte("12345678901234567890123456789012") // 32 bytes
}

func TestEncrypt_Decrypt_RoundTrip(t *testing.T) {
	key := mustKey()
	plain := []byte("secret data")
	ciphertext, err := Encrypt(plain, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(ciphertext) <= len(plain) {
		t.Errorf("ciphertext too short")
	}
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(decrypted, plain) {
		t.Errorf("Decrypt: got %q, want %q", decrypted, plain)
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	_, err := Encrypt([]byte("x"), []byte("short"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrKeyLength) {
		t.Errorf("expected ErrKeyLength, got %v", err)
	}
}

func TestDecrypt_InvalidKeyLength(t *testing.T) {
	_, err := Decrypt([]byte("x"), []byte("short"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrKeyLength) {
		t.Errorf("expected ErrKeyLength, got %v", err)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key := mustKey()
	ciphertext, _ := Encrypt([]byte("data"), key)
	otherKey := []byte("abcdefghijklmnopqrstuvwxyzabcdef") // 32 bytes, different
	_, err := Decrypt(ciphertext, otherKey)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecrypt) {
		t.Errorf("expected ErrDecrypt, got %v", err)
	}
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	key := mustKey()
	_, err := Decrypt([]byte("too short"), key)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecrypt) {
		t.Errorf("expected ErrDecrypt, got %v", err)
	}
}
