package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ulibaysya/casenya/internal/modules/itemtag"
	"github.com/ulibaysya/casenya/pkg/db"
)

func (s *Storage) CreateTag(ctx context.Context, name string, isPublic bool, values []string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}
	defer db.RollbackTX(ctx, tx, s.logger)

	var keyID int64
	const keyInsertSQL = `
		INSERT INTO item_tag_keys (name, is_public) 
		VALUES ($1, $2)
		RETURNING id;
	`
	if err = tx.QueryRow(ctx, keyInsertSQL, name, isPublic).Scan(&keyID); err != nil {
		if isAlreadyExistsErr(err) {
			return itemtag.ErrKeyAlreadyExists
		}
		return fmt.Errorf("inserting tag key: %w", err)
	}

	valuesInsertSQL, valuesInsertArgs := buildValuesInsertSQL(keyID, values)
	if _, err := tx.Exec(ctx, valuesInsertSQL, valuesInsertArgs...); err != nil {
		return fmt.Errorf("inserting tag values: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting tx: %w", err)
	}

	return nil
}

func isAlreadyExistsErr(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.UniqueViolation {
		return false
	}
	const expectedConstraintName = `item_tag_keys_name_key`
	return pgErr.ConstraintName == expectedConstraintName
}

func buildValuesInsertSQL(keyID int64, values []string) (string, []any) {
	b := &strings.Builder{}

	b.WriteString(`
		INSERT INTO item_tag_values (key_id, name)
		VALUES ($1, $2)
	`)
	args := []any{keyID, values[0]}

	argIndex := 3
	for _, value := range values[1:] {
		fmt.Fprintf(b, ", ($1, $%d)", argIndex)
		args = append(args, value)
		argIndex += 1
	}

	b.WriteString(";")

	return b.String(), args
}
