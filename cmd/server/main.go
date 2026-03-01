package main

import (
	"github.com/MaxRadzey/gog/internal/app"
	"github.com/MaxRadzey/gog/internal/config"
)

func main() {
	AppConfig := config.New()

	config.ParseEnv(AppConfig)
	config.ParseFlags(AppConfig)

	if err := app.Run(AppConfig); err != nil {
		panic(err)
	}
}
