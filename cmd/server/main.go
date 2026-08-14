package main

import (
	"log/slog"
	"os"

	"github.com/fark-tee/fark-tee-backend/cmd/server/di"
)

func main() {
	server, cleanup, err := di.Initialize()
	if err != nil {
		slog.Error("Failed to initialize server", slog.Any("error", err))

		os.Exit(1)
	}

	defer cleanup()

	server.Start()
}
