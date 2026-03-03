package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MaxRadzey/gog/internal/client/http_client"
)

type mockSecretGetter struct {
	getFunc func(ctx context.Context, id int64) (*http_client.SecretResponse, error)
}

func (m *mockSecretGetter) GetSecret(ctx context.Context, id int64) (*http_client.SecretResponse, error) {
	return m.getFunc(ctx, id)
}

func TestSecretGetCommand_Execute_Success(t *testing.T) {
	client := &mockSecretGetter{
		getFunc: func(ctx context.Context, id int64) (*http_client.SecretResponse, error) {
			return &http_client.SecretResponse{ID: 1, SecretType: "login", Data: []byte(`{"login":"u"}`)}, nil
		},
	}
	cmd := NewSecretGetCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1"})
	require.NoError(t, err)
	assert.Contains(t, got, "id=1")
	assert.Contains(t, got, "type=login")
}

func TestSecretGetCommand_Execute_NoArgs(t *testing.T) {
	client := &mockSecretGetter{getFunc: func(ctx context.Context, id int64) (*http_client.SecretResponse, error) { return nil, nil }}
	cmd := NewSecretGetCommand(client)
	_, err := cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "args:")
}

func TestSecretGetCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("not found")
	client := &mockSecretGetter{getFunc: func(ctx context.Context, id int64) (*http_client.SecretResponse, error) {
		return nil, clientErr
	}}
	cmd := NewSecretGetCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
