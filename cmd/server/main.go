package main

import (
	"fmt"
	"os"

	"github.com/MaxRadzey/gog/internal/server/app"
	"github.com/MaxRadzey/gog/internal/server/config"
)

func main() {
	cfg := config.New()
	config.ParseEnv(cfg)
	config.ParseFlags(cfg)

	a, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init app: %v\n", err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to run app", err)
		os.Exit(1)
	}
}
