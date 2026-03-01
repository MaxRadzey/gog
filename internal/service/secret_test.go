package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/MaxRadzey/gog/internal/crypto"
	"github.com/MaxRadzey/gog/internal/repository"
	"github.com/MaxRadzey/gog/internal/repository/mocks"
)

var testEncryptionKey = []byte("dev-encryption-key-32bytes-long!")

func TestSecretService_Create(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	payload := []byte(`{"login":"u","password":"p"}`)
	repo.EXPECT().
		Create(gomock.Any(), int64(1), SecretTypeLoginPassword, gomock.Any()).
		Return(int64(1), nil)

	_, err := svc.Create(ctx, 1, SecretTypeLoginPassword, payload)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
}

func TestSecretService_Create_ValidationError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	_, err := svc.Create(ctx, 1, SecretTypeLoginPassword, []byte(`{}`))
	if err == nil {
		t.Fatal("expected validation error")
	}
	var val *ErrValidation
	if !errors.As(err, &val) {
		t.Fatalf("expected *ErrValidation, got %T", err)
	}

	_, err = svc.Create(ctx, 1, "unknown_type", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !errors.As(err, &val) {
		t.Fatalf("expected *ErrValidation, got %T", err)
	}
}

func TestSecretService_GetByID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	plain := []byte(`{"login":"a","password":"b"}`)
	encrypted, _ := crypto.Encrypt(plain, testEncryptionKey)
	secret := &repository.Secret{
		ID: 1, UserID: 1, SecretType: SecretTypeLoginPassword, Data: encrypted,
		CreatedAt: time.Now(), UpdatedAt: time.Now(), IsDeleted: false,
	}

	repo.EXPECT().GetByID(ctx, int64(1), int64(1)).Return(secret, nil)

	dto, err := svc.GetByID(ctx, 1, 1)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if dto.ID != 1 || dto.SecretType != SecretTypeLoginPassword {
		t.Errorf("dto: id=%d type=%s", dto.ID, dto.SecretType)
	}
	if string(dto.Data) != string(plain) {
		t.Errorf("decrypted data = %q, want %q", dto.Data, plain)
	}
}

func TestSecretService_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	repo.EXPECT().GetByID(ctx, int64(999), int64(1)).Return(nil, nil)

	_, err := svc.GetByID(ctx, 999, 1)
	if !errors.Is(err, ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestSecretService_ListByUserID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	plain := []byte(`{"content":"note"}`)
	encrypted, _ := crypto.Encrypt(plain, testEncryptionKey)
	list := []*repository.Secret{
		{ID: 1, UserID: 1, SecretType: SecretTypeText, Data: encrypted, CreatedAt: time.Now(), UpdatedAt: time.Now(), IsDeleted: false},
	}

	repo.EXPECT().ListByUserID(ctx, int64(1)).Return(list, nil)

	dtos, err := svc.ListByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(dtos) != 1 || string(dtos[0].Data) != string(plain) {
		t.Errorf("list len=%d data=%q", len(dtos), dtos[0].Data)
	}
}

func TestSecretService_Update(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	plain := []byte(`{"login":"old","password":"old"}`)
	encrypted, _ := crypto.Encrypt(plain, testEncryptionKey)
	existing := &repository.Secret{ID: 1, UserID: 1, SecretType: SecretTypeLoginPassword, Data: encrypted, CreatedAt: time.Now(), UpdatedAt: time.Now(), IsDeleted: false}
	newPayload := []byte(`{"login":"new","password":"new"}`)

	repo.EXPECT().GetByID(ctx, int64(1), int64(1)).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(nil)

	err := svc.Update(ctx, 1, 1, newPayload)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestSecretService_Delete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockSecretRepository(ctrl)
	svc := NewSecretService(repo, testEncryptionKey)

	repo.EXPECT().Delete(ctx, int64(1), int64(1)).Return(nil)

	err := svc.Delete(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
