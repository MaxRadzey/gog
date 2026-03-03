package config

import (
	"errors"
	"testing"

	"github.com/MaxRadzey/gog/internal/server/crypto"
)

const validEncryptionKey = "dev-encryption-key-32bytes-long!" // 32 bytes

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "postgres://gog:gog@localhost:5433/gog_test",
		EncryptionKey: validEncryptionKey,
		EnableHTTPS:   false,
	}
	if err := Validate(cfg); err != nil {
		t.Errorf("Validate(valid config): got %v, want nil", err)
	}
}

func TestValidate_ValidHTTPS(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "postgres://localhost/db",
		EncryptionKey: validEncryptionKey,
		EnableHTTPS:   true,
		TLSCertFile:   "/path/to/cert.pem",
		TLSKeyFile:    "/path/to/key.pem",
	}
	if err := Validate(cfg); err != nil {
		t.Errorf("Validate(valid HTTPS config): got %v, want nil", err)
	}
}

func TestValidate_EmptyDSN(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "",
		EncryptionKey: validEncryptionKey,
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate(empty DSN): expected error")
	}
	if err.Error() != "config: DATABASE_DSN is required" {
		t.Errorf("Validate(empty DSN): got %q", err.Error())
	}
}

func TestValidate_InvalidEncryptionKey(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "postgres://localhost/db",
		EncryptionKey: "short",
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate(short key): expected error")
	}
	if !errors.Is(err, crypto.ErrKeyLength) {
		t.Errorf("Validate(short key): got %v, want ErrKeyLength", err)
	}
}

func TestValidate_HTTPSWithoutCert(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "postgres://localhost/db",
		EncryptionKey: validEncryptionKey,
		EnableHTTPS:   true,
		TLSCertFile:   "",
		TLSKeyFile:    "/path/to/key.pem",
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate(HTTPS without cert): expected error")
	}
	if err.Error() != "config: TLS_CERT_FILE is required when ENABLE_HTTPS is true" {
		t.Errorf("Validate(HTTPS without cert): got %q", err.Error())
	}
}

func TestValidate_HTTPSWithoutKey(t *testing.T) {
	cfg := &Config{
		DatabaseDSN:   "postgres://localhost/db",
		EncryptionKey: validEncryptionKey,
		EnableHTTPS:   true,
		TLSCertFile:   "/path/to/cert.pem",
		TLSKeyFile:    "",
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate(HTTPS without key): expected error")
	}
	if err.Error() != "config: TLS_KEY_FILE is required when ENABLE_HTTPS is true" {
		t.Errorf("Validate(HTTPS without key): got %q", err.Error())
	}
}
