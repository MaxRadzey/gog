package repository

import (
	"context"
	"time"
)

// User — сущность пользователя из БД.
type User struct {
	ID           int64
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsDeleted    bool
}

// UserRepository — интерфейс репозитория пользователей.
type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (int64, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}
