package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ulibaysya/casenya/internal/modules/user"
)

func (s *Storage) CreateUser(ctx context.Context, name string) (int64, error) {
	var id int64
	const userInsertSQL = `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id;
	`
	err := s.db.QueryRow(ctx, userInsertSQL, name).Scan(&id)
	if err != nil {
		if isAlreadyExistsError(err) {
			return 0, user.ErrAlreadyExists
		}
		return 0, fmt.Errorf("inserting user: %w", err)
	}

	return id, nil
}

func isAlreadyExistsError(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.UniqueViolation {
		return false
	}
	const expectedConstraintName = `users_name_key`
	return pgErr.ConstraintName == expectedConstraintName
}
