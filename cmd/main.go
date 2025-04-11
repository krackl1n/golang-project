package main

import (
	"log/slog"
	"os"

	"github.com/krackl1n/golang-project/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1)
	}
}
