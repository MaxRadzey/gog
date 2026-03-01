package postgres_test

import (
	"context"
	"testing"
)

// mustCreateUser создаёт пользователя и возвращает его id. При ошибке завершает тест.
func mustCreateUser(t *testing.T, repos *testRepos, login string) int64 {
	t.Helper()
	id, err := repos.User.Create(context.Background(), login, "hash")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return id
}

func TestSecretRepository_Create(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	id, err := repos.Secret.Create(ctx, userID, "login_password", []byte("encrypted-payload"))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero ID")
	}

	secret, err := repos.Secret.GetByID(ctx, id, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret == nil {
		t.Fatal("expected secret after Create")
	}
	if secret.ID != id {
		t.Errorf("ID: got %d", secret.ID)
	}
	if secret.UserID != userID {
		t.Errorf("UserID: got %d", secret.UserID)
	}
	if secret.SecretType != "login_password" {
		t.Errorf("SecretType: got %q", secret.SecretType)
	}
	if string(secret.Data) != "encrypted-payload" {
		t.Errorf("Data: got %q", secret.Data)
	}
	if secret.IsDeleted {
		t.Error("expected IsDeleted false")
	}
	if secret.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if secret.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestSecretRepository_Create_EmptyData(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "bob")

	id, err := repos.Secret.Create(ctx, userID, "text", []byte{})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero ID")
	}

	secret, _ := repos.Secret.GetByID(ctx, id, userID)
	if secret == nil {
		t.Fatal("expected secret")
	}
	if len(secret.Data) != 0 {
		t.Errorf("expected empty Data, got len %d", len(secret.Data))
	}
}

func TestSecretRepository_GetByID_Found(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	createdID, err := repos.Secret.Create(ctx, userID, "bank_card", []byte("card-data"))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	secret, err := repos.Secret.GetByID(ctx, createdID, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret == nil {
		t.Fatal("expected secret")
	}
	if secret.ID != createdID {
		t.Errorf("ID: got %d", secret.ID)
	}
	if secret.UserID != userID {
		t.Errorf("UserID: got %d", secret.UserID)
	}
	if secret.SecretType != "bank_card" {
		t.Errorf("SecretType: got %q", secret.SecretType)
	}
	if string(secret.Data) != "card-data" {
		t.Errorf("Data: got %q", secret.Data)
	}
}

func TestSecretRepository_GetByID_NotFound(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	secret, err := repos.Secret.GetByID(ctx, 999, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret != nil {
		t.Errorf("expected nil secret, got %+v", secret)
	}
}

func TestSecretRepository_GetByID_AfterDelete(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	secretID, err := repos.Secret.Create(ctx, userID, "text", []byte("x"))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	err = repos.Secret.Delete(ctx, secretID, userID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	secret, err := repos.Secret.GetByID(ctx, secretID, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret != nil {
		t.Errorf("expected nil for deleted secret, got %+v", secret)
	}
}

func TestSecretRepository_ListByUserID_Empty(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	list, err := repos.Secret.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}

func TestSecretRepository_ListByUserID_Multiple(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	_, _ = repos.Secret.Create(ctx, userID, "login_password", []byte("a"))
	_, _ = repos.Secret.Create(ctx, userID, "text", []byte("b"))
	_, _ = repos.Secret.Create(ctx, userID, "binary", []byte("c"))

	list, err := repos.Secret.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}
	if list[0].SecretType != "login_password" || string(list[0].Data) != "a" {
		t.Errorf("first item: got type %q data %q", list[0].SecretType, list[0].Data)
	}
	if list[1].SecretType != "text" || string(list[1].Data) != "b" {
		t.Errorf("second item: got type %q data %q", list[1].SecretType, list[1].Data)
	}
	if list[2].SecretType != "binary" || string(list[2].Data) != "c" {
		t.Errorf("third item: got type %q data %q", list[2].SecretType, list[2].Data)
	}
}

func TestSecretRepository_ListByUserID_OnlyOwn(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID1 := mustCreateUser(t, repos, "alice")
	userID2 := mustCreateUser(t, repos, "bob")

	_, _ = repos.Secret.Create(ctx, userID1, "text", []byte("alice-secret"))
	_, _ = repos.Secret.Create(ctx, userID2, "text", []byte("bob-secret"))

	list, err := repos.Secret.ListByUserID(ctx, userID1)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item for alice, got %d", len(list))
	}
	if string(list[0].Data) != "alice-secret" {
		t.Errorf("expected alice-secret, got %q", list[0].Data)
	}
}

func TestSecretRepository_Update(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	secretID, err := repos.Secret.Create(ctx, userID, "text", []byte("old"))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = repos.Secret.Update(ctx, secretID, userID, []byte("new"))
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	secret, err := repos.Secret.GetByID(ctx, secretID, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret == nil {
		t.Fatal("expected secret")
	}
	if string(secret.Data) != "new" {
		t.Errorf("Data: got %q", secret.Data)
	}
}

func TestSecretRepository_Delete(t *testing.T) {
	repos := setupDB(t)
	ctx := context.Background()
	userID := mustCreateUser(t, repos, "alice")

	secretID, err := repos.Secret.Create(ctx, userID, "text", []byte("x"))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = repos.Secret.Delete(ctx, secretID, userID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	secret, err := repos.Secret.GetByID(ctx, secretID, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if secret != nil {
		t.Errorf("expected nil after Delete, got %+v", secret)
	}
}
