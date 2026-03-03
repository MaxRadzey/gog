package command

import "context"

// SessionEnder — интерфейс завершения сессии на сервере.
type SessionEnder interface {
	Logout(ctx context.Context) error
}

// LogoutCommand — команда выхода из аккаунта.
type LogoutCommand struct {
	client SessionEnder
}

// NewLogoutCommand создаёт команду логаута.
func NewLogoutCommand(client SessionEnder) *LogoutCommand {
	return &LogoutCommand{
		client: client,
	}
}

// Execute выполняет выход из аккаунта (аргументы не требуются).
func (c *LogoutCommand) Execute(ctx context.Context, args []string) (string, error) {
	err := c.client.Logout(ctx)
	if err != nil {
		return "", err
	}

	return "logout successful", nil
}
