package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Storage struct {
	db     DB
	logger *slog.Logger
}

func New(db DB, logger *slog.Logger) *Storage {
	return &Storage{
		db:     db,
		logger: logger,
	}
}
