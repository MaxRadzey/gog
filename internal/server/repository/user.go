// Package repository — интерфейсы и сущности для хранения пользователей и секретов.
package repository

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mocks/mock_user_repository.go -package=mocks github.com/MaxRadzey/gog/internal/server/repository UserRepository

import (
	"context"
	"time"
)

// User — запись пользователя (логин, хэш пароля, метаданные).
type User struct {
	ID           int64
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsDeleted    bool
}

// UserRepository — создание пользователя, получение по логину и по id.
type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (int64, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}
