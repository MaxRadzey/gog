package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MaxRadzey/gog/internal/server/repository"
	"github.com/MaxRadzey/gog/internal/server/repository/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// StorageOption — функция для настройки Storage.
type StorageOption func(*storageConfig)

type storageConfig struct {
	migrationsPath string
	maxOpenConns   int
	maxIdleConns   int
}

// WithMigrationsPath устанавливает путь к миграциям.
func WithMigrationsPath(path string) StorageOption {
	return func(c *storageConfig) {
		c.migrationsPath = path
	}
}

// WithMaxOpenConns устанавливает максимальное количество открытых соединений.
func WithMaxOpenConns(n int) StorageOption {
	return func(c *storageConfig) {
		c.maxOpenConns = n
	}
}

// WithMaxIdleConns устанавливает максимальное количество простаивающих соединений.
func WithMaxIdleConns(n int) StorageOption {
	return func(c *storageConfig) {
		c.maxIdleConns = n
	}
}

// Storage — подключение к БД и репозитории.
type Storage struct {
	db               *sql.DB
	userRepository   repository.UserRepository
	secretRepository repository.SecretRepository
}

// UserRepository возвращает репозиторий пользователей.
func (s *Storage) UserRepository() repository.UserRepository {
	return s.userRepository
}

// SecretRepository возвращает репозиторий секретов.
func (s *Storage) SecretRepository() repository.SecretRepository {
	return s.secretRepository
}

// InitializeStorage запускает миграции, подключается к БД и создаёт репозитории.
func InitializeStorage(dsn string, opts ...StorageOption) (*Storage, error) {
	if dsn == "" {
		return nil, errors.New("database DSN required")
	}

	cfg := &storageConfig{
		migrationsPath: "",
		maxOpenConns:   0,
		maxIdleConns:   0,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	if err := RunMigrations(dsn, cfg.migrationsPath); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if cfg.maxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.maxOpenConns)
	}
	if cfg.maxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.maxIdleConns)
	}

	return &Storage{
		db:               db,
		userRepository:   postgres.NewUserRepository(db),
		secretRepository: postgres.NewSecretRepository(db),
	}, nil
}

// Close закрывает подключение к БД.
func (s *Storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
