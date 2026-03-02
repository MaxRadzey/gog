// Package handler — общий контейнер для HTTP-хендлеров: сервисы и секрет для подписи куки.
package handler

import (
	"github.com/MaxRadzey/gog/internal/server/service"
)

// Handler передаётся в хендлеры user и secret; хранит сервисы и cookie secret.
type Handler struct {
	Services     *service.Services
	CookieSecret string
}

// New создаёт Handler с переданными сервисами и секретом для куки.
func New(services *service.Services, cookieSecret string) *Handler {
	return &Handler{
		Services:     services,
		CookieSecret: cookieSecret,
	}
}
