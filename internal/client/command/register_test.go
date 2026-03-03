package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRegistrar struct {
	registerFunc func(ctx context.Context, login, password string) error
}

func (m *mockUserRegistrar) Register(ctx context.Context, login, password string) error {
	return m.registerFunc(ctx, login, password)
}

func TestRegisterCommand_Execute_Success(t *testing.T) {
	client := &mockUserRegistrar{
		registerFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
	}

	cmd := NewRegisterCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"user", "pass"})
	require.NoError(t, err)
	assert.Equal(t, "registration successful", got)
}

func TestRegisterCommand_Execute_NotEnoughArgs(t *testing.T) {
	client := &mockUserRegistrar{
		registerFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
	}

	cmd := NewRegisterCommand(client)

	got, err := cmd.Execute(context.Background(), []string{"user"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Contains(t, err.Error(), "args: <login> <password>")

	got, err = cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Empty(t, got)
}

func TestRegisterCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("login already taken")
	client := &mockUserRegistrar{
		registerFunc: func(ctx context.Context, login, password string) error {
			return clientErr
		},
	}

	cmd := NewRegisterCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"user", "pass"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
