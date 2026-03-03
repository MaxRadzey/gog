package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand_Execute(t *testing.T) {
	cmd := NewHelpCommand()
	got, err := cmd.Execute(context.Background(), nil)
	require.NoError(t, err)
	assert.Contains(t, got, "GOG")
	assert.Contains(t, got, "register")
	assert.Contains(t, got, "secret-list")
	assert.Contains(t, got, "help")
}
