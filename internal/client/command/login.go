package command

import (
	"context"
	"errors"
)

// UserAuthenticator — интерфейс аутентификации пользователя на сервере.
type UserAuthenticator interface {
	Login(ctx context.Context, login, password string) error
}

// LoginCommand — команда входа в аккаунт.
type LoginCommand struct {
	client UserAuthenticator
}

// NewLoginCommand создаёт команду логина.
func NewLoginCommand(client UserAuthenticator) *LoginCommand {
	return &LoginCommand{
		client: client,
	}
}

// Execute выполняет вход: args[0] — login, args[1] — password.
func (c *LoginCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("args: <login> <password>")
	}

	err := c.client.Login(ctx, args[0], args[1])
	if err != nil {
		return "", err
	}

	return "login successful", nil
}
