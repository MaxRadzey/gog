package command

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MaxRadzey/gog/internal/client/http_client"
)

// SecretGetter — интерфейс получения секрета по id.
type SecretGetter interface {
	GetSecret(ctx context.Context, id int64) (*http_client.SecretResponse, error)
}

// SecretGetCommand — команда получения секрета по id.
type SecretGetCommand struct {
	client SecretGetter
}

// NewSecretGetCommand создаёт команду получения секрета.
func NewSecretGetCommand(client SecretGetter) *SecretGetCommand {
	return &SecretGetCommand{client: client}
}

// Execute возвращает секрет по id. args[0] — id.
func (c *SecretGetCommand) Execute(ctx context.Context, args []string) (string, error) {
	id, err := parseID(args, "id")
	if err != nil {
		return "", err
	}
	secret, err := c.client.GetSecret(ctx, id)
	if err != nil {
		return "", err
	}
	return formatSecret(secret), nil
}

func formatSecret(s *http_client.SecretResponse) string {
	return fmt.Sprintf("id=%d type=%s created=%s", s.ID, s.SecretType, s.CreatedAt)
}

// parseID парсит id из args[0].
func parseID(args []string, argName string) (int64, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("args: <%s>", argName)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", args[0], err)
	}
	return id, nil
}
