package main

import (
	"log/slog"
	"os"

	"github.com/ryanprayoga/messhub/backend/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		slog.Error("bootstrap app failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer application.Close()

	if err := application.Listen(); err != nil {
		slog.Error("listen failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
