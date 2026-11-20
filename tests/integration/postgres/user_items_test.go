package postgres

import (
	cryptorand "crypto/rand"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/internal/modules/user/storage"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func TestListUserItems_NotFound(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	gotItems, err := storage.New(db, pkgslog.DiscardLogger).ListUserItems(ctx, 12345)
	if !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("got error while listing user items: %v; expected %v", err, user.ErrNotFound)
	}

	if gotItems != nil {
		t.Fatalf("got user items: %v; expected nil", gotItems)
	}
}

func TestListUserItems_Empty(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	var id int64
	err := db.QueryRow(ctx, "INSERT INTO users(name) VALUES ($1) RETURNING id", "some-name").Scan(&id)
	if err != nil {
		t.Fatalf("error while inserting user: %v", err)
	}

	gotItems, err := storage.New(db, pkgslog.DiscardLogger).ListUserItems(ctx, id)
	if err != nil {
		t.Fatalf("got error while listing user items: %v; expected nil", err)
	}

	if gotItems != nil {
		t.Fatalf("got user items: %v; expected nil", gotItems)
	}
}

func TestListUserItems(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	var userID int64
	err := db.QueryRow(ctx, "INSERT INTO users(name) VALUES ($1) RETURNING id", "some-name").Scan(&userID)
	if err != nil {
		t.Fatalf("error while inserting user: %v", err)
	}

	items := make(map[user.Item]struct{}, 5)
	statuses := map[int64]string{
		1: "active",
		2: "sold",
		3: "in_withdrawal",
		4: "withdrawn",
	}
	for range 5 {
		var itemID int64
		if err = db.QueryRow(ctx, "INSERT INTO items(name) VALUES ($1) RETURNING id", cryptorand.Text()).Scan(&itemID); err != nil {
			t.Fatalf("error while inserting item: %v", err)
		}
		price := rand.Int64N(500000) + 1
		status := rand.Int64N(4) + 1
		var userItemID int64
		if err = db.QueryRow(ctx, "INSERT INTO users_items(user_id, item_id, attached_price, status_id) VALUES ($1, $2, $3, $4) RETURNING id", userID, itemID, price, status).Scan(&userItemID); err != nil {
			t.Fatalf("error while inserting user item: %v", err)
		}
		item := user.Item{
			ID:     userItemID,
			ItemID: itemID,
			Price:  price,
			Status: statuses[status],
		}
		items[item] = struct{}{}
	}

	gotItems, err := storage.New(db, pkgslog.DiscardLogger).ListUserItems(ctx, userID)
	if err != nil {
		t.Fatalf("got error while listing user items: %v; expected nil", err)
	}

	if gotItems == nil {
		t.Fatalf("got nil user items, expected non-nil")
	}

	for _, gotItem := range gotItems {
		if _, found := items[gotItem]; !found {
			t.Fatalf("not found user item: %v; expected to be found", gotItem)
		}
		delete(items, gotItem)
	}
	if len(items) != 0 {
		t.Fatalf("not all items were found")
	}
}
