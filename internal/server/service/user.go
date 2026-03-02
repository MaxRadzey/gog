package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/MaxRadzey/gog/internal/server/repository"
)

// UserService — регистрация (с хэшированием пароля) и аутентификация по логину/паролю.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService создаёт сервис с переданным репозиторием.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

const (
	MinLoginLength    = 3
	MinPasswordLength = 6
)

// Register создаёт пользователя (пароль хэшируется bcrypt).
func (s *UserService) Register(ctx context.Context, login, plainPassword string) (userID int64, err error) {
	if login == "" {
		return 0, &ErrValidation{Msg: "login is required"}
	}
	if len(login) < MinLoginLength {
		return 0, &ErrValidation{Msg: "login must be at least " + strconv.Itoa(MinLoginLength) + " characters"}
	}
	if plainPassword == "" {
		return 0, &ErrValidation{Msg: "password is required"}
	}
	if len(plainPassword) < MinPasswordLength {
		return 0, &ErrValidation{Msg: "password must be at least " + strconv.Itoa(MinPasswordLength) + " characters"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, login, string(hash))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, &ErrDuplicateLogin{Login: login}
		}
		return 0, err
	}
	return id, nil
}

// Authenticate проверяет логин и пароль; при успехе возвращает пользователя.
func (s *UserService) Authenticate(ctx context.Context, login, plainPassword string) (*repository.User, error) {
	u, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, &ErrInvalidCredentials{Login: login}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword)); err != nil {
		return nil, &ErrInvalidCredentials{Login: login}
	}
	return u, nil
}
