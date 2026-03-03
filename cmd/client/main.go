package main

import (
	"fmt"
	"os"

	"github.com/MaxRadzey/gog/internal/client/app"
	"github.com/MaxRadzey/gog/internal/client/config"
)

func main() {
	cfg := config.New()
	a, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init app: %v\n", err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
