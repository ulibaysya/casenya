package slog

import (
	"context"
	"log/slog"
)

func Error(ctx context.Context, logger *slog.Logger, msg string, err error) {
	logger.LogAttrs(ctx, slog.LevelError, msg, Err("error", err))
}
