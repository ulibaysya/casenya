package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/pkg/db"
)

func (s *Storage) ListUserItems(ctx context.Context, id int64) ([]user.Item, error) {
	const itemsSelectSQL = `
		SELECT
			u.id AS user_id,
			i.id,
			i.item_id,
			i.attached_price,
			s.code
		FROM users AS u
		LEFT JOIN users_items AS i ON i.user_id = u.id
		LEFT JOIN users_items_statuses AS s ON i.status_id = s.id
		WHERE u.id = $1;
	`
	rows, err := s.db.Query(ctx, itemsSelectSQL, id)
	if err != nil {
		return nil, fmt.Errorf("selecting user items: %w", err)
	}
	type itemsSelectRow struct {
		UserID        int64
		UserItemID    pgtype.Int8
		ItemID        pgtype.Int8
		AttachedPrice pgtype.Int8
		StatusCode    pgtype.Text
	}
	itemsSelectRows, err := pgx.CollectRows(rows, pgx.RowToStructByPos[itemsSelectRow])
	if err != nil {
		return nil, fmt.Errorf("collecting user items rows: %w", err)
	}

	if len(itemsSelectRows) == 0 {
		return nil, user.ErrNotFound
	}

	if !itemsSelectRows[0].UserItemID.Valid {
		return nil, nil
	}

	items := make([]user.Item, 0, len(itemsSelectRows))
	for _, row := range itemsSelectRows {
		items = append(items, user.Item{
			ID:     row.UserItemID.Int64,
			ItemID: row.ItemID.Int64,
			Price:  row.AttachedPrice.Int64,
			Status: row.StatusCode.String,
		})
	}

	return items, nil
}

func (s *Storage) SellUserItem(ctx context.Context, userID int64, itemID int64) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("beginning transaction: %w", err)
	}
	defer db.RollbackTX(ctx, tx, s.logger)

	const checkUserSQL = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
			FOR KEY SHARE
		);
	`
	var userExists bool
	if err = tx.QueryRow(ctx, checkUserSQL, userID).Scan(&userExists); err != nil {
		return 0, fmt.Errorf("selecting user: %w", err)
	}
	if !userExists {
		return 0, user.ErrNotFound
	}

	const checkItemSQL = `
		SELECT 
			s.code
		FROM users_items AS i
		JOIN users_items_statuses AS s ON s.id = i.status_id
		WHERE i.id = $1
			AND i.user_id = $2
		FOR UPDATE OF i;
	`
	var code string
	if err = tx.QueryRow(ctx, checkItemSQL, itemID, userID).Scan(&code); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, user.ErrItemNotFound
		}
		return 0, fmt.Errorf("selecting user item: %w", err)
	}
	if user.ItemStatus(code) != user.ItemStatusActive {
		return 0, user.ErrItemNotActive
	}

	const statusUpdateSQL = `
		UPDATE users_items
		SET status_id = (
			SELECT s.id
			FROM users_items_statuses AS s
			WHERE s.code = 'sold'
		)
		WHERE id = $1
			AND user_id = $2
			AND status_id = (
				SELECT s.id
				FROM users_items_statuses AS s
				WHERE s.code = 'active'
			)
		RETURNING user_id, attached_price;
	`
	var actualUserID int64
	var attachedPrice int64
	if err = tx.QueryRow(ctx, statusUpdateSQL, itemID, userID).Scan(&actualUserID, &attachedPrice); err != nil {
		return 0, fmt.Errorf("updating user item status: %w", err)
	}

	const balanceUpdateSQL = `
		UPDATE users
		set balance = balance + $1
		where id = $2
		RETURNING balance;
	`
	var balance int64
	if err = tx.QueryRow(ctx, balanceUpdateSQL, attachedPrice, actualUserID).Scan(&balance); err != nil {
		return 0, fmt.Errorf("updating user balance: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commiting transaction: %w", err)
	}

	return balance, nil
}
