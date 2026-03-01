package postgres

import (
	"context"
	"database/sql"

	"github.com/MaxRadzey/gog/internal/repository"
)

// UserRepository — репозиторий пользователей в PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт репозиторий пользователей.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create сохраняет пользователя и возвращает его id.
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

// GetByLogin возвращает пользователя по логину (только не удалённого).
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

// GetByID возвращает пользователя по id (только не удалённого).
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
