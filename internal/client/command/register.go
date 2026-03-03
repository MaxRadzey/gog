package command

import (
	"context"
	"errors"
)

// UserRegistrar — интерфейс регистрации пользователя на сервере.
type UserRegistrar interface {
	Register(ctx context.Context, login, password string) error
}

// RegisterCommand — команда регистрации пользователя.
type RegisterCommand struct {
	client UserRegistrar
}

// NewRegisterCommand создаёт команду регистрации.
func NewRegisterCommand(client UserRegistrar) *RegisterCommand {
	return &RegisterCommand{
		client: client,
	}
}

// Execute выполняет регистрацию: args[0] — login, args[1] — password.
func (c *RegisterCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("args: <login> <password>")
	}

	err := c.client.Register(ctx, args[0], args[1])
	if err != nil {
		return "", err
	}

	return "registration successful", nil
}
