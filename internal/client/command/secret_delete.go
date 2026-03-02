package command

import (
	"context"
	"fmt"
)

// SecretDeleter — интерфейс удаления секрета.
type SecretDeleter interface {
	DeleteSecret(ctx context.Context, id int64) error
}

// SecretDeleteCommand — команда удаления секрета.
type SecretDeleteCommand struct {
	client SecretDeleter
}

// NewSecretDeleteCommand создаёт команду удаления секрета.
func NewSecretDeleteCommand(client SecretDeleter) *SecretDeleteCommand {
	return &SecretDeleteCommand{client: client}
}

// Execute удаляет секрет. args[0] — id.
func (c *SecretDeleteCommand) Execute(ctx context.Context, args []string) (string, error) {
	id, err := parseID(args, "id")
	if err != nil {
		return "", err
	}
	err = c.client.DeleteSecret(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("deleted secret id=%d", id), nil
}
