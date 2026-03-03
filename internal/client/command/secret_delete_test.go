package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSecretDeleter struct {
	deleteFunc func(ctx context.Context, id int64) error
}

func (m *mockSecretDeleter) DeleteSecret(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

func TestSecretDeleteCommand_Execute_Success(t *testing.T) {
	client := &mockSecretDeleter{
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	cmd := NewSecretDeleteCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1"})
	require.NoError(t, err)
	assert.Equal(t, "deleted secret id=1", got)
}

func TestSecretDeleteCommand_Execute_NoArgs(t *testing.T) {
	client := &mockSecretDeleter{deleteFunc: func(ctx context.Context, id int64) error { return nil }}
	cmd := NewSecretDeleteCommand(client)
	_, err := cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "args:")
}

func TestSecretDeleteCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("not found")
	client := &mockSecretDeleter{deleteFunc: func(ctx context.Context, id int64) error {
		return clientErr
	}}
	cmd := NewSecretDeleteCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
