package secret

import "encoding/json"

// CreateRequest — тело запроса создания секрета.
type CreateRequest struct {
	SecretType string          `json:"secret_type" binding:"required"`
	Data       json.RawMessage `json:"data" binding:"required"`
}

// CreateResponse — ответ после создания секрета.
type CreateResponse struct {
	ID int64 `json:"id"`
}

// UpdateRequest — тело запроса обновления секрета (только data).
type UpdateRequest struct {
	Data json.RawMessage `json:"data" binding:"required"`
}

// SecretResponse — один секрет в ответе (data — расшифрованный JSON).
type SecretResponse struct {
	ID         int64           `json:"id"`
	UserID     int64           `json:"user_id"`
	SecretType string          `json:"secret_type"`
	Data       json.RawMessage `json:"data"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}
