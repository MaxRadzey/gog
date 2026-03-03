package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserAuthenticator struct {
	loginFunc func(ctx context.Context, login, password string) error
}

func (m *mockUserAuthenticator) Login(ctx context.Context, login, password string) error {
	return m.loginFunc(ctx, login, password)
}

func TestLoginCommand_Execute_Success(t *testing.T) {
	client := &mockUserAuthenticator{
		loginFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
	}

	cmd := NewLoginCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"user", "pass"})
	require.NoError(t, err)
	assert.Equal(t, "login successful", got)
}

func TestLoginCommand_Execute_NotEnoughArgs(t *testing.T) {
	client := &mockUserAuthenticator{
		loginFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
	}

	cmd := NewLoginCommand(client)

	got, err := cmd.Execute(context.Background(), []string{"user"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Contains(t, err.Error(), "args: <login> <password>")

	got, err = cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Empty(t, got)
}

func TestLoginCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("invalid credentials")
	client := &mockUserAuthenticator{
		loginFunc: func(ctx context.Context, login, password string) error {
			return clientErr
		},
	}

	cmd := NewLoginCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"user", "pass"})
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
