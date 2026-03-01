package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/MaxRadzey/gog/internal/repository"
	"github.com/MaxRadzey/gog/internal/repository/mocks"
)

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			Create(gomock.Any(), "alice", gomock.Any()).
			Return(int64(1), nil)

		id, err := svc.Register(ctx, "alice", "password123")
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		if id != 1 {
			t.Errorf("id = %d, want 1", id)
		}
	})

	t.Run("duplicate login", func(t *testing.T) {
		repo.EXPECT().
			Create(gomock.Any(), "taken", gomock.Any()).
			Return(int64(0), &pgconn.PgError{Code: "23505"})

		id, err := svc.Register(ctx, "taken", "pass")
		if err == nil {
			t.Fatal("expected error")
		}
		if id != 0 {
			t.Errorf("id = %d, want 0", id)
		}
		var dup *ErrDuplicateLogin
		if !errors.As(err, &dup) {
			t.Fatalf("expected *ErrDuplicateLogin, got %T", err)
		}
		if dup.Login != "taken" {
			t.Errorf("ErrDuplicateLogin.Login = %q, want taken", dup.Login)
		}
	})
}

func TestUserService_Authenticate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			GetByLogin(gomock.Any(), "alice").
			Return(&repository.User{
				ID:           1,
				Login:        "alice",
				PasswordHash: string(hash),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				IsDeleted:    false,
			}, nil)

		u, err := svc.Authenticate(ctx, "alice", "secret")
		if err != nil {
			t.Fatalf("Authenticate: %v", err)
		}
		if u == nil {
			t.Fatal("expected user")
		}
		if u.Login != "alice" || u.ID != 1 {
			t.Errorf("user: id=%d login=%s", u.ID, u.Login)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().
			GetByLogin(gomock.Any(), "nobody").
			Return(nil, nil)

		u, err := svc.Authenticate(ctx, "nobody", "any")
		if err == nil {
			t.Fatal("expected error")
		}
		if u != nil {
			t.Error("expected nil user")
		}
		var creds *ErrInvalidCredentials
		if !errors.As(err, &creds) {
			t.Fatalf("expected *ErrInvalidCredentials, got %T", err)
		}
		if creds.Login != "nobody" {
			t.Errorf("ErrInvalidCredentials.Login = %q, want nobody", creds.Login)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		repo.EXPECT().
			GetByLogin(gomock.Any(), "alice").
			Return(&repository.User{
				ID:           1,
				Login:        "alice",
				PasswordHash: string(hash),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				IsDeleted:    false,
			}, nil)

		u, err := svc.Authenticate(ctx, "alice", "wrongpassword")
		if err == nil {
			t.Fatal("expected error")
		}
		if u != nil {
			t.Error("expected nil user")
		}
		var creds *ErrInvalidCredentials
		if !errors.As(err, &creds) {
			t.Fatalf("expected *ErrInvalidCredentials, got %T", err)
		}
		if creds.Login != "alice" {
			t.Errorf("ErrInvalidCredentials.Login = %q, want alice", creds.Login)
		}
	})
}
