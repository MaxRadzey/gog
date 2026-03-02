// Package postgres — реализация репозиториев пользователей и секретов для PostgreSQL.
package postgres

import (
	"context"
	"database/sql"

	"github.com/MaxRadzey/gog/internal/server/repository"
)

// UserRepository хранит и читает пользователей в таблице users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт репозиторий по переданному *sql.DB.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create вставляет пользователя (login, password_hash), возвращает id.
func (r *UserRepository) Create(ctx context.Context, login, passwordHash string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`,
		login, passwordHash,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByLogin возвращает пользователя по логину; не удалённого.
func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*repository.User, error) {
	var u repository.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, login, password_hash, created_at, updated_at, is_deleted
		 FROM users WHERE login = $1 AND NOT is_deleted`,

		login,
	).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByID возвращает пользователя по id; не удалённого.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*repository.User, error) {
	var u repository.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, login, password_hash, created_at, updated_at, is_deleted
		 FROM users WHERE id = $1 AND NOT is_deleted`,
		id,
	).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
