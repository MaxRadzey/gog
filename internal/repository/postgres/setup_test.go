package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/MaxRadzey/gog/internal/repository/postgres"
	"github.com/MaxRadzey/gog/internal/storage"
)

var testDB *sql.DB

// testRepos — все репозитории для тестов. Один setup, один TRUNCATE.
type testRepos struct {
	User *postgres.UserRepository
	DB   *sql.DB
}

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://gog:gog@localhost:5433/gog_test?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		os.Exit(1)
	}

	migrationsPath := os.Getenv("TEST_MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath, _ = filepath.Abs("../../../migrations")
	}
	if err := storage.RunMigrations(dsn, migrationsPath); err != nil {
		_ = db.Close()
		os.Exit(1)
	}

	testDB = db
	code := m.Run()
	_ = db.Close()
	os.Exit(code)
}

// setupDB очищает таблицу users и возвращает репозитории для тестов.
func setupDB(t *testing.T) *testRepos {
	t.Helper()
	_, err := testDB.ExecContext(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return &testRepos{
		User: postgres.NewUserRepository(testDB),
		DB:   testDB,
	}
}
