package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ulibaysya/casenya/internal/modules/user"
)

func (s *Storage) SaveUserOpaqueToken(ctx context.Context, id int64, token string) error {
	const tokenInsertSQL = `
		INSERT INTO users_opaque_tokens (user_id, token)
		VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, tokenInsertSQL, id, token)
	if err != nil {
		if isNotFoundErr(err) {
			return user.ErrNotFound
		}
		return fmt.Errorf("inserting user opaque token: %w", err)
	}

	return nil
}

func isNotFoundErr(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.ForeignKeyViolation {
		return false
	}
	const expectedConstraintName = `users_opaque_tokens_user_id_fkey`
	return pgErr.ConstraintName == expectedConstraintName
}

func (s *Storage) ExchangeUserOpaqueToken(ctx context.Context, token string) (int64, error) {
	const tokenDeleteSQl = `
		DELETE FROM users_opaque_tokens
		WHERE token = $1
		RETURNING user_id;
	`
	var userID int64
	err := s.db.QueryRow(ctx, tokenDeleteSQl, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, user.ErrTokenNotFound
		}
		return 0, fmt.Errorf("deleting token: %w", err)
	}

	return userID, nil
}
