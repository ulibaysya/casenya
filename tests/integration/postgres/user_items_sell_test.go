package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/internal/modules/user/storage"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func TestSellUserItem_NotFound(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).SellUserItem(ctx, 123, 456)
	if !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("got error: %v; expected %v", err, user.ErrNotFound)
	}

	if gotBalance != 0 {
		t.Fatalf("got balance: %v; expected 0", gotBalance)
	}
}

func TestSellUserItem_ItemNotFound(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	balance := int64(123456)
	userID := db.addUser("1", "1", balance)

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).SellUserItem(ctx, userID, 456)
	if !errors.Is(err, user.ErrItemNotFound) {
		t.Fatalf("got error: %v; expected %v", err, user.ErrItemNotFound)
	}
	if gotBalance != 0 {
		t.Fatalf("got balance: %v; expected 0", gotBalance)
	}

	var balanceAfterSell int64
	if err = db.QueryRow(ctx, "SELECT balance FROM users").Scan(&balanceAfterSell); err != nil {
		t.Fatalf("error while querying balance: %v", err)
	}
	if balanceAfterSell != balance {
		t.Fatalf("got balance after sell: %v; expected: %v", balanceAfterSell, balance)
	}
}

func TestSellUserItem_ItemNotActive(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	balance := int64(123583190)
	userID := db.addUser("1", "1", balance)
	itemID := db.addItem("1")
	attachedPrice := int64(47128)
	status := user.ItemStatusWithdrawn
	userItemID := db.addUserItem(userID, itemID, attachedPrice, status)

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).SellUserItem(ctx, userID, userItemID)
	if !errors.Is(err, user.ErrItemNotActive) {
		t.Fatalf("got error: %v; expected %v", err, user.ErrItemNotActive)
	}
	if gotBalance != 0 {
		t.Fatalf("got balance: %v; expected 0", gotBalance)
	}

	var balanceAfterSell int64
	if err = db.QueryRow(ctx, "SELECT balance FROM users").Scan(&balanceAfterSell); err != nil {
		t.Fatalf("error while querying balance: %v", err)
	}
	if balanceAfterSell != balance {
		t.Fatalf("got balance after sell: %v; expected: %v", balanceAfterSell, balance)
	}

	var statusAfterSell string
	if err = db.QueryRow(ctx, "SELECT s.code FROM users_items AS i JOIN users_items_statuses AS s ON i.status_id = s.id").Scan(&statusAfterSell); err != nil {
		t.Fatalf("error while querying status: %v", err)
	}
	if user.ItemStatus(statusAfterSell) != status {
		t.Fatalf("got status after sell: %v; expected: %v", statusAfterSell, status)
	}
}

func TestSellUserItem(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	balance := int64(123583190)
	userID := db.addUser("1", "1", balance)
	itemID := db.addItem("1")
	attachedPrice := int64(47128)
	status := user.ItemStatusActive
	userItemID := db.addUserItem(userID, itemID, attachedPrice, status)

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).SellUserItem(ctx, userID, userItemID)
	if err != nil {
		t.Fatalf("got error: %v; expected nil", err)
	}
	if gotBalance != balance+attachedPrice {
		t.Fatalf("got balance: %v; expected %v", gotBalance, balance+attachedPrice)
	}

	var balanceAfterSell int64
	if err = db.QueryRow(ctx, "SELECT balance FROM users").Scan(&balanceAfterSell); err != nil {
		t.Fatalf("error while querying balance: %v", err)
	}
	expectedPrice := balance + attachedPrice
	if balanceAfterSell != expectedPrice {
		t.Fatalf("got balance after sell: %v; expected: %v", balanceAfterSell, expectedPrice)
	}

	var statusAfterSell string
	if err = db.QueryRow(ctx, "SELECT s.code FROM users_items AS i JOIN users_items_statuses AS s ON i.status_id = s.id").Scan(&statusAfterSell); err != nil {
		t.Fatalf("error while querying status: %v", err)
	}
	if user.ItemStatus(statusAfterSell) != user.ItemStatusSold {
		t.Fatalf("got status after sell: %v; expected: %v", statusAfterSell, user.ItemStatusSold)
	}
}

func TestSellUserItem_AlienItem(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	balance := int64(123456)
	userOwningItemID := db.addUser("1", "1", balance)
	userSellingItemID := db.addUser("2", "2", balance)

	itemID := db.addItem("1")
	attachedPrice := int64(47128)
	status := user.ItemStatusActive
	userItemID := db.addUserItem(userOwningItemID, itemID, attachedPrice, status)

	gotBalance, err := storage.New(db, pkgslog.DiscardLogger).SellUserItem(ctx, userSellingItemID, userItemID)
	if !errors.Is(err, user.ErrItemNotFound) {
		t.Fatalf("got error: %v; expected %v", err, user.ErrItemNotFound)
	}
	if gotBalance != 0 {
		t.Fatalf("got balance: %v; expected 0", gotBalance)
	}

	rows, err := db.Query(ctx, "SELECT balance FROM users")
	if err != nil {
		t.Fatalf("error while querying balance: %v", err)
	}
	balances, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		t.Fatalf("error while collecting balances: %v", err)
	}
	if len(balances) != 2 {
		t.Fatalf("got %v balances, expected 2", len(balances))
	}
	if balances[0] != balances[1] || balances[0] != balance {
		t.Fatalf("got balances: %v and %v; expected both %v", balances[0], balances[1], balance)
	}

	var statusAfterSell string
	if err = db.QueryRow(ctx, "SELECT s.code FROM users_items AS i JOIN users_items_statuses AS s ON i.status_id = s.id").Scan(&statusAfterSell); err != nil {
		t.Fatalf("error while querying status: %v", err)
	}
	if user.ItemStatus(statusAfterSell) != status {
		t.Fatalf("got status after sell: %v; expected: %v", statusAfterSell, status)
	}
}
