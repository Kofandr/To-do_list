package logger

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		opts.Level = slog.LevelDebug
	case "WARN":
		opts.Level = slog.LevelWarn
	case "ERROR":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
