package http_client

import "encoding/json"

// RegisterRequest — тело запроса POST /api/user/register.
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest — тело запроса POST /api/user/login.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateSecretRequest — тело запроса POST /api/secret.
type CreateSecretRequest struct {
	SecretType string          `json:"secret_type"`
	Data       json.RawMessage `json:"data"`
}

// CreateSecretResponse — ответ после создания секрета.
type CreateSecretResponse struct {
	ID int64 `json:"id"`
}

// UpdateSecretRequest — тело запроса PUT /api/secret/:id.
type UpdateSecretRequest struct {
	Data json.RawMessage `json:"data"`
}

// SecretResponse — один секрет в ответе (GET /api/secret, GET /api/secret/:id).
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
