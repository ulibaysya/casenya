package postgres

import (
	"errors"
	"testing"

	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/internal/modules/user/storage"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func TestChangeUserBalance_InsufficientBalance(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	var id int64
	err := db.QueryRow(ctx, `INSERT INTO users(name, balance) VALUES ('some-name', 100) RETURNING id`).Scan(&id)
	if err != nil {
		t.Fatalf("error while inserting user: %v", err)
	}

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).ChangeUserBalance(ctx, id, -1000)
	if !errors.Is(err, user.ErrInsufficientBalance) {
		t.Fatalf("got error while changing user balance: %v; expected %v", err, user.ErrInsufficientBalance)
	}

	if gotBalance != 0 {
		t.Fatalf("got balance %v, expected 0", gotBalance)
	}
}

func TestChangeUserBalance_NotFound(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).ChangeUserBalance(ctx, 123, 456789)
	if !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("got error while changing user balance: %v; expected %v", err, user.ErrNotFound)
	}

	if gotBalance != 0 {
		t.Fatalf("got balance %v, expected 0", gotBalance)
	}
}
