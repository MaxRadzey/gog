package service

import (
	"errors"
	"fmt"
)

// ErrDuplicateLogin — логин уже занят при регистрации.
type ErrDuplicateLogin struct {
	Login string
}

func (e *ErrDuplicateLogin) Error() string {
	return fmt.Sprintf("duplicate login: %s", e.Login)
}

// ErrInvalidCredentials — неверный логин или пароль при аутентификации.
type ErrInvalidCredentials struct {
	Login string
}

func (e *ErrInvalidCredentials) Error() string {
	return fmt.Sprintf("invalid credentials: %s", e.Login)
}

// ErrValidation — ошибка валидации входных данных (пустые поля, длина и т.д.).
type ErrValidation struct {
	Msg string
}

func (e *ErrValidation) Error() string {
	return e.Msg
}

// ErrSecretNotFound — секрет не найден или не принадлежит пользователю.
var ErrSecretNotFound = errors.New("secret not found")
