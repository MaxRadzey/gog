package service

import (
	"errors"
	"fmt"
)

// ErrDuplicateLogin — при регистрации логин уже существует.
type ErrDuplicateLogin struct {
	Login string
}

func (e *ErrDuplicateLogin) Error() string {
	return fmt.Sprintf("duplicate login: %s", e.Login)
}

// ErrInvalidCredentials — неверный логин или пароль при входе.
type ErrInvalidCredentials struct {
	Login string
}

func (e *ErrInvalidCredentials) Error() string {
	return fmt.Sprintf("invalid credentials: %s", e.Login)
}

// ErrValidation — ошибка валидации (формат полей, ограничения).
type ErrValidation struct {
	Msg string
}

func (e *ErrValidation) Error() string {
	return e.Msg
}

// ErrSecretNotFound — секрет не найден по id или не принадлежит пользователю.
var ErrSecretNotFound = errors.New("secret not found")
