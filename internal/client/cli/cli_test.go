package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/MaxRadzey/gog/internal/client/command"
)

func TestRun_ExecuteCommand(t *testing.T) {
	registry := command.CommandRegistry{
		"help": command.NewHelpCommand(),
	}
	stdinR, stdinW, _ := os.Pipe()
	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()

	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdin = stdinR
	os.Stdout = stdoutW
	os.Stderr = stderrW
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	_, _ = stdinW.Write([]byte("help\n"))
	_ = stdinW.Close()

	Run(context.Background(), registry)

	_ = stdoutW.Close()
	_ = stderrW.Close()
	out, _ := io.ReadAll(stdoutR)
	errOut, _ := io.ReadAll(stderrR)

	if !strings.Contains(string(out), ">") {
		t.Errorf("stdout: expected prompt >, got %q", out)
	}
	if !strings.Contains(string(out), "register") {
		t.Errorf("stdout: expected help text with 'register', got %q", out)
	}
	if len(errOut) != 0 {
		t.Errorf("stderr: expected empty, got %q", errOut)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	registry := command.CommandRegistry{
		"help": command.NewHelpCommand(),
	}
	stdinR, stdinW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()

	oldStdin := os.Stdin
	oldStderr := os.Stderr
	os.Stdin = stdinR
	os.Stderr = stderrW
	defer func() {
		os.Stdin = oldStdin
		os.Stderr = oldStderr
	}()

	_, _ = stdinW.Write([]byte("unknowncmd\n"))
	_ = stdinW.Close()

	done := make(chan struct{})
	go func() {
		Run(context.Background(), registry)
		_ = stderrW.Close()
		close(done)
	}()

	stderrBuf := &bytes.Buffer{}
	_, _ = io.Copy(stderrBuf, stderrR)
	<-done

	errStr := stderrBuf.String()
	if !strings.Contains(errStr, "Unknown command") {
		t.Errorf("stderr: expected 'Unknown command', got %q", errStr)
	}
}
