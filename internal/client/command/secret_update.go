package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// SecretUpdater — интерфейс обновления секрета.
type SecretUpdater interface {
	UpdateSecret(ctx context.Context, id int64, data json.RawMessage) error
}

// SecretUpdateCommand — команда обновления секрета.
type SecretUpdateCommand struct {
	client SecretUpdater
}

// NewSecretUpdateCommand создаёт команду обновления секрета.
func NewSecretUpdateCommand(client SecretUpdater) *SecretUpdateCommand {
	return &SecretUpdateCommand{client: client}
}

// Execute обновляет секрет. args[0] — id, далее пары key=value для полей.
func (c *SecretUpdateCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("args: <id> key=value [key=value ...]")
	}
	id, err := parseID(args[:1], "id")
	if err != nil {
		return "", err
	}
	data, err := parseKeyValueArgs(args[1:])
	if err != nil {
		return "", err
	}
	err = c.client.UpdateSecret(ctx, id, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("updated secret id=%d", id), nil
}
