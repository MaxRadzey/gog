package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MaxRadzey/gog/internal/server/crypto"
	"github.com/MaxRadzey/gog/internal/server/repository"
)

// SecretDTO — секрет с расшифрованными данными для отдачи клиенту.
type SecretDTO struct {
	ID         int64
	UserID     int64
	SecretType string
	Data       []byte
	CreatedAt  string
	UpdatedAt  string
}

// SecretService — бизнес-логика работы с секретами.
type SecretService struct {
	repo repository.SecretRepository
	key  []byte
}

// NewSecretService создаёт сервис секретов.
func NewSecretService(repo repository.SecretRepository, key []byte) *SecretService {
	return &SecretService{repo: repo, key: key}
}

// Create создаёт секрет: валидирует payload по типу, шифрует и сохраняет.
func (s *SecretService) Create(ctx context.Context, userID int64, secretType string, payloadJSON []byte) (int64, error) {
	if err := ValidatePayloadByType(secretType, payloadJSON); err != nil {
		return 0, err
	}
	encrypted, err := crypto.Encrypt(payloadJSON, s.key)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, userID, secretType, encrypted)
}

// GetByID возвращает секрет с расшифрованными данными.
func (s *SecretService) GetByID(ctx context.Context, id, userID int64) (*SecretDTO, error) {
	secret, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrSecretNotFound
	}
	decrypted, err := crypto.Decrypt(secret.Data, s.key)
	if err != nil {
		return nil, err
	}
	return secretToDTO(secret, decrypted), nil
}

// ListByUserID возвращает все секреты пользователя с расшифрованными данными.
func (s *SecretService) ListByUserID(ctx context.Context, userID int64) ([]*SecretDTO, error) {
	list, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*SecretDTO, 0, len(list))
	for _, sec := range list {
		decrypted, err := crypto.Decrypt(sec.Data, s.key)
		if err != nil {
			return nil, err
		}
		out = append(out, secretToDTO(sec, decrypted))
	}
	return out, nil
}

// Update обновляет данные секрета.
func (s *SecretService) Update(ctx context.Context, id, userID int64, payloadJSON []byte) error {
	secret, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if secret == nil {
		return ErrSecretNotFound
	}
	if err := ValidatePayloadByType(secret.SecretType, payloadJSON); err != nil {
		return err
	}
	encrypted, err := crypto.Encrypt(payloadJSON, s.key)
	if err != nil {
		return err
	}
	err = s.repo.Update(ctx, id, userID, encrypted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSecretNotFound
		}
		return err
	}
	return nil
}

// Delete помечает секрет как удалённый.
func (s *SecretService) Delete(ctx context.Context, id, userID int64) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSecretNotFound
		}
		return err
	}
	return nil
}

func secretToDTO(s *repository.Secret, decrypted []byte) *SecretDTO {
	return &SecretDTO{
		ID:         s.ID,
		UserID:     s.UserID,
		SecretType: s.SecretType,
		Data:       decrypted,
		CreatedAt:  s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
