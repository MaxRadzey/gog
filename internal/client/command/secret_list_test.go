package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MaxRadzey/gog/internal/client/http_client"
)

type mockSecretLister struct {
	listFunc func(ctx context.Context) ([]http_client.SecretResponse, error)
}

func (m *mockSecretLister) ListSecrets(ctx context.Context) ([]http_client.SecretResponse, error) {
	return m.listFunc(ctx)
}

func TestSecretListCommand_Execute_Success(t *testing.T) {
	client := &mockSecretLister{
		listFunc: func(ctx context.Context) ([]http_client.SecretResponse, error) {
			return []http_client.SecretResponse{
				{ID: 1, SecretType: "login", CreatedAt: "2025-01-01"},
			}, nil
		},
	}
	cmd := NewSecretListCommand(client)
	got, err := cmd.Execute(context.Background(), nil)
	require.NoError(t, err)
	assert.Contains(t, got, "id=1")
	assert.Contains(t, got, "type=login")
}

func TestSecretListCommand_Execute_EmptyList(t *testing.T) {
	client := &mockSecretLister{
		listFunc: func(ctx context.Context) ([]http_client.SecretResponse, error) {
			return nil, nil
		},
	}
	cmd := NewSecretListCommand(client)
	got, err := cmd.Execute(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, "no secrets", got)
}

func TestSecretListCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("unauthorized")
	client := &mockSecretLister{listFunc: func(ctx context.Context) ([]http_client.SecretResponse, error) {
		return nil, clientErr
	}}
	cmd := NewSecretListCommand(client)
	got, err := cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
