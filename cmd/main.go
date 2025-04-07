package main

import (
	"log/slog"
	"os"

	"github.com/krackl1n/golang-project/internal/app"
)

func main() {
	// Настройка логера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if err := app.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
