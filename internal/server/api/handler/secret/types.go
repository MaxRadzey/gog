package secret

import "encoding/json"

// CreateRequest — запрос на создание секрета (тип + data в JSON).
type CreateRequest struct {
	SecretType string          `json:"secret_type" binding:"required"`
	Data       json.RawMessage `json:"data" binding:"required" swaggertype:"object"`
}

// CreateResponse — ответ создания, возвращает id секрета.
type CreateResponse struct {
	ID int64 `json:"id"`
}

// UpdateRequest — обновление секрета, только поле data.
type UpdateRequest struct {
	Data json.RawMessage `json:"data" binding:"required" swaggertype:"object"`
}

// SecretResponse — один секрет в ответе API (data — JSON payload).
type SecretResponse struct {
	ID         int64           `json:"id"`
	UserID     int64           `json:"user_id"`
	SecretType string          `json:"secret_type"`
	Data       json.RawMessage `json:"data" swaggertype:"object"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}
