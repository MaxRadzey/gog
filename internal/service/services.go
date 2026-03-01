package service

import "github.com/MaxRadzey/gog/internal/repository"

// Services — контейнер сервисов приложения.
type Services struct {
	User *UserService
}

// ServicesOption — функция для настройки Services.
type ServicesOption func(*Services)

// WithUserService переопределяет сервис пользователей.
func WithUserService(svc *UserService) ServicesOption {
	return func(s *Services) {
		s.User = svc
	}
}

// NewServices создаёт все сервисы из контейнера репозиториев.
// Сервисы зависят от интерфейса RepositoryContainer, а не от конкретной реализации Storage.
func NewServices(container repository.RepositoryContainer, opts ...ServicesOption) *Services {
	s := &Services{
		User: NewUserService(container.UserRepository()),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
