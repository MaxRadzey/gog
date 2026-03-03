// Package app — точка входа клиента: инициализация логгера и HTTP-клиента, сборка реестра команд и запуск интерактивного CLI.
package app

import (
	"context"
	"time"

	"github.com/MaxRadzey/gog/internal/client/cli"
	"github.com/MaxRadzey/gog/internal/client/command"
	"github.com/MaxRadzey/gog/internal/client/config"
	"github.com/MaxRadzey/gog/internal/client/http_client"
	"github.com/MaxRadzey/gog/internal/logger"
)

// App хранит конфиг и HTTP-клиент; реестр команд собирается в Run.
type App struct {
	cfg    *config.Config
	client *http_client.Client
}

// New инициализирует логгер и HTTP-клиент по конфигу, возвращает приложение.
func New(cfg *config.Config) (*App, error) {
	if err := logger.Initialize("info"); err != nil {
		return nil, err
	}
	client, err := http_client.New(cfg.ServerURL, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, client: client}, nil
}

// Run запускает цикл CLI (чтение команд из stdin и выполнение).
func (a *App) Run() error {
	registry := command.CommandRegistry{
		"help":          command.NewHelpCommand(),
		"version":       command.NewVersionCommand(),
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
