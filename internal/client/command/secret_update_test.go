package command

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSecretUpdater struct {
	updateFunc func(ctx context.Context, id int64, data json.RawMessage) error
}

func (m *mockSecretUpdater) UpdateSecret(ctx context.Context, id int64, data json.RawMessage) error {
	return m.updateFunc(ctx, id, data)
}

func TestSecretUpdateCommand_Execute_Success(t *testing.T) {
	client := &mockSecretUpdater{
		updateFunc: func(ctx context.Context, id int64, data json.RawMessage) error {
			return nil
		},
	}
	cmd := NewSecretUpdateCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1", "login=newuser"})
	require.NoError(t, err)
	assert.Equal(t, "updated secret id=1", got)
}

func TestSecretUpdateCommand_Execute_NotEnoughArgs(t *testing.T) {
	client := &mockSecretUpdater{updateFunc: func(ctx context.Context, id int64, data json.RawMessage) error { return nil }}
	cmd := NewSecretUpdateCommand(client)
	_, err := cmd.Execute(context.Background(), []string{"1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "args:")
}

func TestSecretUpdateCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("forbidden")
	client := &mockSecretUpdater{updateFunc: func(ctx context.Context, id int64, data json.RawMessage) error {
		return clientErr
	}}
	cmd := NewSecretUpdateCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"1", "login=x"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
