package command

import (
	"context"
	"fmt"

	"github.com/MaxRadzey/gog/internal/client/http_client"
)

// SecretLister — интерфейс получения списка секретов.
type SecretLister interface {
	ListSecrets(ctx context.Context) ([]http_client.SecretResponse, error)
}

// SecretListCommand — команда вывода списка секретов.
type SecretListCommand struct {
	client SecretLister
}

// NewSecretListCommand создаёт команду списка секретов.
func NewSecretListCommand(client SecretLister) *SecretListCommand {
	return &SecretListCommand{client: client}
}

// Execute выводит список секретов пользователя. Аргументы не требуются.
func (c *SecretListCommand) Execute(ctx context.Context, args []string) (string, error) {
	list, err := c.client.ListSecrets(ctx)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "no secrets", nil
	}
	var out string
	for i, s := range list {
		if i > 0 {
			out += "\n"
		}
		out += fmt.Sprintf("id=%d type=%s created=%s", s.ID, s.SecretType, s.CreatedAt)
	}
	return out, nil
}
