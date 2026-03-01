package service

import "fmt"

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
