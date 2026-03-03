package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSessionEnder struct {
	logoutFunc func(ctx context.Context) error
}

func (m *mockSessionEnder) Logout(ctx context.Context) error {
	return m.logoutFunc(ctx)
}

func TestLogoutCommand_Execute_Success(t *testing.T) {
	client := &mockSessionEnder{
		logoutFunc: func(ctx context.Context) error {
			return nil
		},
	}

	cmd := NewLogoutCommand(client)
	got, err := cmd.Execute(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, "logout successful", got)
}

func TestLogoutCommand_Execute_Success_IgnoresArgs(t *testing.T) {
	client := &mockSessionEnder{
		logoutFunc: func(ctx context.Context) error {
			return nil
		},
	}

	cmd := NewLogoutCommand(client)
	got, err := cmd.Execute(context.Background(), []string{"extra", "args"})
	require.NoError(t, err)
	assert.Equal(t, "logout successful", got)
}

func TestLogoutCommand_Execute_ClientError(t *testing.T) {
	clientErr := errors.New("session error")
	client := &mockSessionEnder{
		logoutFunc: func(ctx context.Context) error {
			return clientErr
		},
	}

	cmd := NewLogoutCommand(client)
	got, err := cmd.Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Empty(t, got)
	assert.Same(t, clientErr, err)
}
