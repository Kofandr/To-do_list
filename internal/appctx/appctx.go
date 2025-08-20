package appctx

import (
	"context"
	"log/slog"
)

type loggerKeyType struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKeyType{}, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	val := ctx.Value(loggerKeyType{})

	logg, ok := val.(*slog.Logger)
	if !ok || logg == nil {
		return slog.Default()
	}

	return logg
}
