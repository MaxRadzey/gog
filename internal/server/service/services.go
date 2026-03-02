package service

import "github.com/MaxRadzey/gog/internal/server/repository"

// Services — контейнер сервисов приложения.
type Services struct {
	User   *UserService
	Secret *SecretService
}

// ServicesOption — функция для настройки Services.
type ServicesOption func(*Services)

// WithUserService переопределяет сервис пользователей.
func WithUserService(svc *UserService) ServicesOption {
	return func(s *Services) {
		s.User = svc
	}
}

// WithSecretService переопределяет сервис секретов.
func WithSecretService(svc *SecretService) ServicesOption {
	return func(s *Services) {
		s.Secret = svc
	}
}

// NewServices создаёт все сервисы из контейнера репозиториев.
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
