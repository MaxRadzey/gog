package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand_Execute(t *testing.T) {
	cmd := NewVersionCommand()

	got, err := cmd.Execute(context.Background(), nil)

	require.NoError(t, err)
	assert.Contains(t, got, "version ")
	assert.Contains(t, got, "build ")
}
