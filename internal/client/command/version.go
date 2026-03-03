package command

import (
	"context"
	"fmt"

	"github.com/MaxRadzey/gog/internal/client/version"
)

// VersionCommand выводит версию и дату сборки клиента.
type VersionCommand struct{}

// NewVersionCommand создаёт команду version.
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

// Execute возвращает строку с версией и датой сборки.
func (c *VersionCommand) Execute(ctx context.Context, args []string) (string, error) {
	return fmt.Sprintf("version %s\nbuild %s", version.Version, version.BuildDate), nil
}
