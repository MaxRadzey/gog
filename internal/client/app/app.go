// Package app собирает конфиг, логгер, HTTP-клиент, реестр команд и запускает CLI.
package app

import (
	"context"

	"github.com/MaxRadzey/gog/internal/client/cli"
	"github.com/MaxRadzey/gog/internal/client/command"
	"github.com/MaxRadzey/gog/internal/client/config"
	"github.com/MaxRadzey/gog/internal/client/http_client"
	"github.com/MaxRadzey/gog/internal/logger"
)

// App — клиентское приложение: конфиг и HTTP-клиент.
type App struct {
	cfg    *config.Config
	client *http_client.Client
}

// New создаёт приложение: инициализирует логгер и HTTP-клиент, собирает реестр команд.
func New(cfg *config.Config) (*App, error) {
	if err := logger.Initialize("info"); err != nil {
		return nil, err
	}
	client, err := http_client.New(cfg.ServerURL)
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, client: client}, nil
}

// Run запускает интерактивный CLI с реестром команд.
func (a *App) Run() error {
	registry := command.CommandRegistry{
		"help":          command.NewHelpCommand(),
		"register":      command.NewRegisterCommand(a.client),
		"login":         command.NewLoginCommand(a.client),
		"logout":        command.NewLogoutCommand(a.client),
		"secret-list":   command.NewSecretListCommand(a.client),
		"secret-get":    command.NewSecretGetCommand(a.client),
		"secret-create": command.NewSecretCreateCommand(a.client),
		"secret-update": command.NewSecretUpdateCommand(a.client),
		"secret-delete": command.NewSecretDeleteCommand(a.client),
	}
	cli.Run(context.Background(), registry)
	return nil
}
