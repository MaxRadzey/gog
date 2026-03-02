package command

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSecretCreator struct {
	createFunc func(ctx context.Context, secretType string, data json.RawMessage) (int64, error)
}

func (m *mockSecretCreator) CreateSecret(ctx context.Context, secretType string, data json.RawMessage) (int64, error) {
	return m.createFunc(ctx, secretType, data)
}

func TestSecretCreateCommand_Execute_Success(t *testing.T) {
	client := &mockSecretCreator{
		createFunc: func(ctx context.Context, secretType string, data json.RawMessage) (int64, error) {
			return 42, nil
		},
	}
	cmd := NewSecretCreateCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"login", "login=u", "password=p"})
	require.NoError(t, err)
	assert.Equal(t, "created secret id=42", got)
}

func TestSecretCreateCommand_Execute_NoType(t *testing.T) {
	client := &mockSecretCreator{createFunc: func(ctx context.Context, secretType string, data json.RawMessage) (int64, error) { return 0, nil }}
	cmd := NewSecretCreateCommand(client)
	_, err := cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "args:")
}

func TestSecretCreateCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("conflict")
	client := &mockSecretCreator{createFunc: func(ctx context.Context, secretType string, data json.RawMessage) (int64, error) {
		return 0, clientErr
	}}
	cmd := NewSecretCreateCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"login", "login=u"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
