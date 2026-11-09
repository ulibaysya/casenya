package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"github.com/ulibaysya/casenya/pkg/db"
)

func (s *Storage) CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("beginning tx: %w", err)
	}
	defer db.RollbackTX(ctx, tx, s.logger)

	var itemID int64
	const itemInsertSQL = `
		INSERT INTO items(name)
		VALUES ($1)
		RETURNING id;
	`
	if err = tx.QueryRow(ctx, itemInsertSQL, name).Scan(&itemID); err != nil {
		if isAlreadyExistsErr(err) {
			return 0, item.ErrAlreadyExists
		}
		return 0, fmt.Errorf("inserting item: %w", err)
	}

	const priceInsertSQL = `
		INSERT INTO items_prices (item_id, price, provider_id)
		SELECT $1, $2, p.id 
		FROM items_price_providers AS p 
		WHERE code = $3;
	`
	res, err := tx.Exec(ctx, priceInsertSQL, itemID, price, item.PPCustom)
	if err != nil {
		return 0, fmt.Errorf("inserting item price: %w", err)
	}
	if aff := res.RowsAffected(); aff != 1 {
		return 0, fmt.Errorf("expected 1 affected row when inserting item price, got %d", aff)
	}

	if len(tags) != 0 {
		tagsInsertSQL, tagsInsertArgs := buildTagsInsertSQL(itemID, tags)
		res, err = tx.Exec(ctx, tagsInsertSQL, tagsInsertArgs...)
		if err != nil {
			return 0, fmt.Errorf("inserting tags: %w", err)
		}
		expectedTagsInsertAff := 0
		for _, tag := range tags {
			expectedTagsInsertAff += len(tag.Values)
		}
		if aff := res.RowsAffected(); aff != int64(expectedTagsInsertAff) {
			return 0, item.ErrTagNotFound
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commiting tx: %w", err)
	}

	return itemID, nil
}

func isAlreadyExistsErr(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.UniqueViolation {
		return false
	}
	const expectedConstraintName = `items_name_key`
	return pgErr.ConstraintName == expectedConstraintName
}

func buildTagsInsertSQL(itemID int64, tags []item.Tag) (string, []any) {
	b := &strings.Builder{}

	b.WriteString(`
		WITH values_ids AS (
			SELECT v.id AS id
			FROM item_tag_values AS v
			JOIN item_tag_keys AS k ON k.id = v.key_id
			JOIN (
				VALUES
					($1, $2)
	`)

	args := []any{tags[0].Key, tags[0].Values[0]}

	latestKeyIndex := 1
	latestValueIndex := latestKeyIndex + 2
	for _, value := range tags[0].Values[1:] {
		fmt.Fprintf(b, `,($%d, $%d)`, latestKeyIndex, latestValueIndex)
		args = append(args, value)
		latestValueIndex += 1
	}
	latestKeyIndex = latestValueIndex

	for _, tag := range tags[1:] {
		args = append(args, tag.Key)
		latestValueIndex := latestKeyIndex + 1
		for _, value := range tag.Values {
			fmt.Fprintf(b, `,($%d, $%d)`, latestKeyIndex, latestValueIndex)
			args = append(args, value)
			latestValueIndex += 1
		}
		latestKeyIndex = latestValueIndex
	}
	itemIDIndex := latestKeyIndex

	fmt.Fprintf(b, `
			) AS n(key_name, value_name) 
			ON n.key_name = k.name AND n.value_name = v.name
		)
		INSERT INTO item_tag_relations
		SELECT $%d, v.id
		FROM values_ids AS v;
	`, itemIDIndex)
	args = append(args, itemID)

	return b.String(), args
}
