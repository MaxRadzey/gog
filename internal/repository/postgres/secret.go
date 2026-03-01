package postgres

import (
	"context"
	"database/sql"

	"github.com/MaxRadzey/gog/internal/repository"
)

// SecretRepository — репозиторий секретов в PostgreSQL.
type SecretRepository struct {
	db *sql.DB
}

// NewSecretRepository создаёт репозиторий секретов.
func NewSecretRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

// Create сохраняет секрет и возвращает его id.
func (r *SecretRepository) Create(ctx context.Context, userID int64, secretType string, data []byte) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO secrets (user_id, secret_type, data) VALUES ($1, $2, $3) RETURNING id`,
		userID, secretType, data,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByID возвращает секрет по id, только если он принадлежит пользователю и не удалён.
func (r *SecretRepository) GetByID(ctx context.Context, id, userID int64) (*repository.Secret, error) {
	var s repository.Secret
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, secret_type, data, created_at, updated_at, is_deleted
		 FROM secrets WHERE id = $1 AND user_id = $2 AND NOT is_deleted`,
		id, userID,
	).Scan(&s.ID, &s.UserID, &s.SecretType, &s.Data, &s.CreatedAt, &s.UpdatedAt, &s.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// ListByUserID возвращает все не удалённые секреты пользователя.
func (r *SecretRepository) ListByUserID(ctx context.Context, userID int64) ([]*repository.Secret, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, secret_type, data, created_at, updated_at, is_deleted
		 FROM secrets WHERE user_id = $1 AND NOT is_deleted ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*repository.Secret
	for rows.Next() {
		var s repository.Secret
		if err := rows.Scan(&s.ID, &s.UserID, &s.SecretType, &s.Data, &s.CreatedAt, &s.UpdatedAt, &s.IsDeleted); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, rows.Err()
}

// Update обновляет данные секрета (только свой и не удалённый).
func (r *SecretRepository) Update(ctx context.Context, id, userID int64, data []byte) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE secrets SET data = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3 AND NOT is_deleted`,
		data, id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete помечает секрет как удалённый.
func (r *SecretRepository) Delete(ctx context.Context, id, userID int64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE secrets SET is_deleted = TRUE, updated_at = NOW() WHERE id = $1 AND user_id = $2 AND NOT is_deleted`,
		id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
