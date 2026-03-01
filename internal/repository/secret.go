package repository

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mocks/mock_secret_repository.go -package=mocks github.com/MaxRadzey/gog/internal/repository SecretRepository

import (
	"context"
	"time"
)

// Secret — сущность секрета из БД (данные в data зашифрованы на уровне приложения).
type Secret struct {
	ID         int64
	UserID     int64
	SecretType string
	Data       []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsDeleted  bool
}

// SecretRepository — интерфейс репозитория секретов.
type SecretRepository interface {
	Create(ctx context.Context, userID int64, secretType string, data []byte) (int64, error)
	GetByID(ctx context.Context, id, userID int64) (*Secret, error)
	ListByUserID(ctx context.Context, userID int64) ([]*Secret, error)
	Update(ctx context.Context, id, userID int64, data []byte) error
	Delete(ctx context.Context, id, userID int64) error
}
