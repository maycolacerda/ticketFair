package services

import (
	"log/slog"
	"os"
)

func Log() {

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)
	slog.Info("Logging Started")
}
