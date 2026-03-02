package postgres_test

import (
	"context"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	id, err := repo.Create(ctx, "alice", "hash123")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero ID")
	}

	user, err := repo.GetByLogin(ctx, "alice")
	if err != nil {
		t.Fatalf("GetByLogin: %v", err)
	}
	if user == nil {
		t.Fatal("expected user after Create")
	}
	if user.ID != id {
		t.Errorf("ID: got %d", user.ID)
	}
	if user.Login != "alice" {
		t.Errorf("Login: got %q", user.Login)
	}
	if user.PasswordHash != "hash123" {
		t.Errorf("PasswordHash: got %q", user.PasswordHash)
	}
	if user.IsDeleted {
		t.Error("expected IsDeleted false")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestUserRepository_Create_DuplicateLogin(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	_, err := repo.Create(ctx, "Max", "hash1")
	if err != nil {
		t.Fatalf("first Create: %v", err)
	}

	_, err = repo.Create(ctx, "Max", "hash2")
	if err == nil {
		t.Fatal("expected error on duplicate login")
	}
}

func TestUserRepository_GetByLogin_Found(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	createdID, err := repo.Create(ctx, "Max", "secret")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	user, err := repo.GetByLogin(ctx, "Max")
	if err != nil {
		t.Fatalf("GetByLogin: %v", err)
	}
	if user == nil {
		t.Fatal("expected user")
	}
	if user.ID != createdID {
		t.Errorf("ID: got %d", user.ID)
	}
	if user.Login != "Max" {
		t.Errorf("Login: got %q", user.Login)
	}
	if user.PasswordHash != "secret" {
		t.Errorf("PasswordHash: got %q", user.PasswordHash)
	}
}

func TestUserRepository_GetByLogin_NotFound(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	user, err := repo.GetByLogin(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetByLogin: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}

func TestUserRepository_GetByID_Found(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	createdID, err := repo.Create(ctx, "bob", "hash")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	user, err := repo.GetByID(ctx, createdID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if user == nil {
		t.Fatal("expected user")
	}
	if user.ID != createdID {
		t.Errorf("ID: got %d", user.ID)
	}
	if user.Login != "bob" {
		t.Errorf("Login: got %q", user.Login)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repos := setupDB(t)
	repo := repos.User
	ctx := context.Background()

	user, err := repo.GetByID(ctx, 999)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}
