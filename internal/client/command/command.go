// Package command — команды CLI: help, register, login, logout, secret-list/get/create/update/delete.
// Каждая команда реализует интерфейс Command и вызывается по имени из реестра.
package command

import "context"

// Command — одна команда: принимает контекст и аргументы, возвращает строку или ошибку.
type Command interface {
	Execute(ctx context.Context, args []string) (string, error)
}

// CommandRegistry — карта имя команды → команда (используется в CLI для вызова по имени).
type CommandRegistry map[string]Command
