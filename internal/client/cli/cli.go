package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MaxRadzey/gog/internal/client/command"
	"github.com/MaxRadzey/gog/internal/logger"
)

func Run(ctx context.Context, registry command.CommandRegistry) {
	reader := bufio.NewScanner(os.Stdin)

	logger.Log.Info("GOG started")
	fmt.Print("> ")

	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if len(line) == 0 {
			fmt.Print("> ")

			continue
		}

		parts := strings.Fields(line)
		cmdName := parts[0]
		args := parts[1:]

		cmd, ok := registry[cmdName]
		if !ok {
			fmt.Fprintln(os.Stderr, "Unknown command:", cmdName)
			fmt.Print("> ")
			continue
		}

		result, err := cmd.Execute(ctx, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			fmt.Print("> ")
			continue
		}

		fmt.Println(result)
		fmt.Print("> ")
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
