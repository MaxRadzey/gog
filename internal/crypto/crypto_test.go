package crypto

import (
	"bytes"
	"testing"
)

var testKey = []byte("dev-encryption-key-32bytes-long!")

func TestEncrypt_Decrypt_RoundTrip(t *testing.T) {
	plain := []byte("secret json payload")
	ciphertext, err := Encrypt(plain, testKey)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if bytes.Equal(ciphertext, plain) {
		t.Error("ciphertext must not equal plaintext")
	}
	if len(ciphertext) <= len(plain) {
		t.Errorf("ciphertext too short (includes nonce)")
	}

	decrypted, err := Decrypt(ciphertext, testKey)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(decrypted, plain) {
		t.Errorf("decrypted %q != plain %q", decrypted, plain)
	}
}

func TestEncrypt_WrongKeyLength(t *testing.T) {
	_, err := Encrypt([]byte("x"), []byte("short"))
	if err != ErrKeyLength {
		t.Errorf("expected ErrKeyLength, got %v", err)
	}
}

func TestDecrypt_WrongKeyLength(t *testing.T) {
	_, err := Decrypt([]byte("x"), []byte("short"))
	if err != ErrKeyLength {
		t.Errorf("expected ErrKeyLength, got %v", err)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	plain := []byte("data")
	ciphertext, err := Encrypt(plain, testKey)
	if err != nil {
		t.Fatal(err)
	}
	wrongKey := []byte("wrong-encryption-key-32bytes-lon")
	_, err = Decrypt(ciphertext, wrongKey)
	if err != ErrDecrypt {
		t.Errorf("expected ErrDecrypt, got %v", err)
	}
}

func TestEncrypt_EmptyPlaintext(t *testing.T) {
	ciphertext, err := Encrypt(nil, testKey)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	decrypted, err := Decrypt(ciphertext, testKey)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if decrypted != nil {
		t.Errorf("expected nil decrypted, got %v", decrypted)
	}
}
