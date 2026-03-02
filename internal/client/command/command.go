package command

import "context"

// Command — команда CLI с методом Execute.
type Command interface {
	Execute(ctx context.Context, args []string) (string, error)
}

// CommandRegistry — реестр команд по имени (имя команды → команда).
type CommandRegistry map[string]Command
