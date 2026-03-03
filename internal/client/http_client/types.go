package http_client

import "encoding/json"

// RegisterRequest — запрос на регистрацию (login + password).
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest — запрос на вход.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateSecretRequest — запрос на создание секрета (тип + данные в JSON).
type CreateSecretRequest struct {
	SecretType string          `json:"secret_type"`
	Data       json.RawMessage `json:"data"`
}

// CreateSecretResponse — ответ создания секрета, возвращает id.
type CreateSecretResponse struct {
	ID int64 `json:"id"`
}

// UpdateSecretRequest — обновление секрета, только поле data.
type UpdateSecretRequest struct {
	Data json.RawMessage `json:"data"`
}

// SecretResponse — один секрет в ответе списка или get по id.
type SecretResponse struct {
	ID         int64           `json:"id"`
	UserID     int64           `json:"user_id"`
	SecretType string          `json:"secret_type"`
	Data       json.RawMessage `json:"data"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

// errResponse — ответ с ошибкой от сервера (4xx/5xx).
type errResponse struct {
	Error string `json:"error"`
}
