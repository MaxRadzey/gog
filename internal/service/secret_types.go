// Типы и обязательные поля для payload секретов (JSON приходит в бизнес-логику, валидируется и шифруется).
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/MaxRadzey/gog/internal/constant"
	"github.com/go-playground/validator/v10"
)

var secretValidator = validator.New()

// expiryRegex — MM/YY (месяц 01–12, год две цифры).
var expiryRegex = regexp.MustCompile(`^(0[1-9]|1[0-2])/\d{2}$`)

func init() {
	_ = secretValidator.RegisterValidation("expiry", func(fl validator.FieldLevel) bool {
		return expiryRegex.MatchString(fl.Field().String())
	})
}

// LoginPasswordPayload — данные для типа login_password. Обязательные: Login, Password.
type LoginPasswordPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Meta     string `json:"meta,omitempty"`
}

// TextPayload — данные для типа text.
type TextPayload struct {
	Content string `json:"content"`
	Meta    string `json:"meta,omitempty"`
}

// BinaryPayload — данные для типа binary. Обязательное: Base64.
type BinaryPayload struct {
	Base64 string `json:"base64" validate:"required"`
	Meta   string `json:"meta,omitempty"`
}

// BankCardPayload — данные для типа bank_card. Number 13–19 цифр, Expiry MM/YY, CVV 3 цифры.
type BankCardPayload struct {
	Number string `json:"number" validate:"required,numeric,min=13,max=19"`
	Holder string `json:"holder" validate:"required"`
	Expiry string `json:"expiry" validate:"required,expiry"`
	CVV    string `json:"cvv,omitempty" validate:"omitempty,len=3,numeric"`
	Meta   string `json:"meta,omitempty"`
}

// validatePayload проверяет структуру по тегам validate. При ошибке возвращает ErrValidation.
func validatePayload(s interface{}) error {
	if s == nil {
		return &ErrValidation{Msg: "payload is required"}
	}
	err := secretValidator.Struct(s)
	if err == nil {
		return nil
	}
	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); ok {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		return &ErrValidation{Msg: strings.Join(msgs, "; ")}
	}
	return &ErrValidation{Msg: err.Error()}
}

// ValidateLoginPasswordPayload проверяет payload по тегам validate. Возвращает ErrValidation при ошибке.
func ValidateLoginPasswordPayload(p *LoginPasswordPayload) error {
	return validatePayload(p)
}

// ValidateTextPayload проверяет payload (nil не допускается).
func ValidateTextPayload(p *TextPayload) error {
	return validatePayload(p)
}

// ValidateBinaryPayload проверяет payload по тегам validate.
func ValidateBinaryPayload(p *BinaryPayload) error {
	return validatePayload(p)
}

// ValidateBankCardPayload проверяет payload по тегам validate.
func ValidateBankCardPayload(p *BankCardPayload) error {
	return validatePayload(p)
}

// UnknownSecretTypeError возвращает ErrValidation для неизвестного типа.
func UnknownSecretTypeError(secretType string) error {
	return &ErrValidation{Msg: "unknown secret type: " + secretType}
}

// ValidatePayloadByType разбирает payloadJSON по secretType, валидирует и возвращает ошибку при неверных данных.
func ValidatePayloadByType(secretType string, payloadJSON []byte) error {
	switch secretType {
	case constant.SecretTypeLoginPassword:
		var p LoginPasswordPayload
		if err := json.Unmarshal(payloadJSON, &p); err != nil {
			return &ErrValidation{Msg: "invalid JSON: " + err.Error()}
		}
		return ValidateLoginPasswordPayload(&p)
	case constant.SecretTypeText:
		var p TextPayload
		if err := json.Unmarshal(payloadJSON, &p); err != nil {
			return &ErrValidation{Msg: "invalid JSON: " + err.Error()}
		}
		return ValidateTextPayload(&p)
	case constant.SecretTypeBinary:
		var p BinaryPayload
		if err := json.Unmarshal(payloadJSON, &p); err != nil {
			return &ErrValidation{Msg: "invalid JSON: " + err.Error()}
		}
		return ValidateBinaryPayload(&p)
	case constant.SecretTypeBankCard:
		var p BankCardPayload
		if err := json.Unmarshal(payloadJSON, &p); err != nil {
			return &ErrValidation{Msg: "invalid JSON: " + err.Error()}
		}
		return ValidateBankCardPayload(&p)
	default:
		return UnknownSecretTypeError(secretType)
	}
}
