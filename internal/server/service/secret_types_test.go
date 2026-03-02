package service

import (
	"testing"
)

func TestValidateLoginPasswordPayload(t *testing.T) {
	if err := ValidateLoginPasswordPayload(nil); err == nil {
		t.Error("expected error for nil")
	}
	if err := ValidateLoginPasswordPayload(&LoginPasswordPayload{}); err == nil {
		t.Error("expected error for empty login")
	}
	if err := ValidateLoginPasswordPayload(&LoginPasswordPayload{Login: "a"}); err == nil {
		t.Error("expected error for empty password")
	}
	if err := ValidateLoginPasswordPayload(&LoginPasswordPayload{Login: "a", Password: "b"}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateTextPayload(t *testing.T) {
	if err := ValidateTextPayload(nil); err == nil {
		t.Error("expected error for nil")
	}
	if err := ValidateTextPayload(&TextPayload{}); err != nil {
		t.Errorf("expected nil for empty content, got %v", err)
	}
}

func TestValidateBinaryPayload(t *testing.T) {
	if err := ValidateBinaryPayload(nil); err == nil {
		t.Error("expected error for nil")
	}
	if err := ValidateBinaryPayload(&BinaryPayload{}); err == nil {
		t.Error("expected error for empty base64")
	}
	if err := ValidateBinaryPayload(&BinaryPayload{Base64: "YQ=="}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateBankCardPayload(t *testing.T) {
	if err := ValidateBankCardPayload(nil); err == nil {
		t.Error("expected error for nil")
	}
	if err := ValidateBankCardPayload(&BankCardPayload{}); err == nil {
		t.Error("expected error for empty fields")
	}
	if err := ValidateBankCardPayload(&BankCardPayload{Number: "1234567890123"}); err == nil {
		t.Error("expected error for missing holder/expiry")
	}
	if err := ValidateBankCardPayload(&BankCardPayload{Number: "1234567890123", Holder: "John", Expiry: "12/25"}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if err := ValidateBankCardPayload(&BankCardPayload{Number: "1234567890123", Holder: "John", Expiry: "13/25"}); err == nil {
		t.Error("expected error for invalid expiry month")
	}
	if err := ValidateBankCardPayload(&BankCardPayload{Number: "1234567890123", Holder: "John", Expiry: "12/25", CVV: "123"}); err != nil {
		t.Errorf("expected nil with valid CVV, got %v", err)
	}
	if err := ValidateBankCardPayload(&BankCardPayload{Number: "1234567890123", Holder: "John", Expiry: "12/25", CVV: "12"}); err == nil {
		t.Error("expected error for CVV not 3 digits")
	}
}
