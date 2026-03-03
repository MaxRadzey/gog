// Package service — бизнес-логика: регистрация и аутентификация пользователей, CRUD секретов с валидацией и шифрованием.
package service

import "github.com/MaxRadzey/gog/internal/server/repository"

// Services — контейнер сервисов User и Secret.
type Services struct {
	User   *UserService
	Secret *SecretService
}

// ServicesOption задаёт опции при создании Services (подмена сервисов в тестах).
type ServicesOption func(*Services)

// WithUserService подставляет свой UserService (например мок).
func WithUserService(svc *UserService) ServicesOption {
	return func(s *Services) {
		s.User = svc
	}
}

// WithSecretService подставляет свой SecretService (например мок).
func WithSecretService(svc *SecretService) ServicesOption {
	return func(s *Services) {
		s.Secret = svc
	}
}

// NewServices создаёт UserService и SecretService по контейнеру репозиториев и ключу шифрования.
func NewServices(container repository.RepositoryContainer, encryptionKey string, opts ...ServicesOption) *Services {
	key := []byte(encryptionKey)
	s := &Services{
		User:   NewUserService(container.UserRepository()),
		Secret: NewSecretService(container.SecretRepository(), key),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
