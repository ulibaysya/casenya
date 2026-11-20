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

func (s *Storage) GetUserBalance(ctx context.Context, id int64) (int64, error) {
	const selectBalanceSQL = `
		SELECT balance
		FROM users
		WHERE user = $1;
	`
	var balance int64
	err := s.db.QueryRow(ctx, selectBalanceSQL, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, user.ErrNotFound
		}
		return 0, fmt.Errorf("selecting user balance: %w", err)
	}

	return balance, nil
}

func (s *Storage) SetUserBalance(ctx context.Context, id int64, balance int64) error {
	const updateBalanceSQL = `
		UPDATE users
		SET balance = $1
		WHERE id = $2;
	`
	res, err := s.db.Exec(ctx, updateBalanceSQL, balance, id)
	if err != nil {
		if isInsufficientBalanceError(err) {
			return user.ErrInsufficientBalance
		}
		return fmt.Errorf("updating user balance: %w", err)
	}
	switch aff := res.RowsAffected(); aff {
	case 0:
		return user.ErrNotFound
	case 1:
	default:
		return fmt.Errorf("unexpected number of affected rows when updating user balance: %v", aff)
	}

	return nil
}

func (s *Storage) ChangeUserBalance(ctx context.Context, id int64, delta int64) (int64, error) {
	const updateBalanceSQL = `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance;
	`
	var balance int64
	err := s.db.QueryRow(ctx, updateBalanceSQL, delta, id).Scan(&balance)
	if err != nil {
		if isInsufficientBalanceError(err) {
			return 0, user.ErrInsufficientBalance
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, user.ErrNotFound
		}
		return 0, fmt.Errorf("updating user balance: %w", err)
	}

	return balance, nil
}

func isInsufficientBalanceError(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.CheckViolation {
		return false
	}
	const expectedConstraintName = `users_balance_check`
	return pgErr.ConstraintName == expectedConstraintName
}
