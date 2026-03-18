// services/logger.go
package services

import (
	"log/slog"
	"os"
)

func InitLogger() {
	env := os.Getenv("GIN_MODE")

	var handler slog.Handler

	if env == "release" {
		// Production — JSON format, parsed by Promtail and indexed in Loki
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Development — human readable, coloured output
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	slog.SetDefault(slog.New(handler))
	slog.Info("Logger initialized", "mode", env)
}
