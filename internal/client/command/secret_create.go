package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// SecretCreator — интерфейс создания секрета.
type SecretCreator interface {
	CreateSecret(ctx context.Context, secretType string, data json.RawMessage) (int64, error)
}

// SecretCreateCommand — команда создания секрета.
type SecretCreateCommand struct {
	client SecretCreator
}

// NewSecretCreateCommand создаёт команду создания секрета.
func NewSecretCreateCommand(client SecretCreator) *SecretCreateCommand {
	return &SecretCreateCommand{client: client}
}

// Execute создаёт секрет. args[0] — secret_type (login, card, text, file и т.д.), далее пары key=value.
func (c *SecretCreateCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("args: <secret_type> [key=value ...]")
	}
	secretType := args[0]
	data, err := parseKeyValueArgs(args[1:])
	if err != nil {
		return "", err
	}
	id, err := c.client.CreateSecret(ctx, secretType, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("created secret id=%d", id), nil
}
