package services

import (
	"log/slog"
	"os"
)

func Log() {
	//remove o timestamp dos logs.
	removeTime := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		return a
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: removeTime,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)
	slog.Info("Logger Initialized...")
}
