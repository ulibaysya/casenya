package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func RollbackTX(ctx context.Context, tx pgx.Tx, logger *slog.Logger) {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		pkgslog.Error(ctx, logger, "error while rollbacking transaction", err)
	}
}
