package postgres

import (
	"errors"
	"testing"

	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/internal/modules/user/storage"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func TestFindUser(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	const someName = "someName"
	const someImageURL = "https://strg.casenya.net/user-images/iHaiAOH6oQGyQPL8IS8WiUpP6M1GION72kJGXmyq0_4.png"
	var id int64
	err := db.QueryRow(ctx, "INSERT INTO users(name, image_url) VALUES ($1, $2) RETURNING id", someName, someImageURL).Scan(&id)
	if err != nil {
		t.Fatalf("error while inserting user: %v", err)
	}

	user, err := storage.New(db, pkgslog.DiscardLogger).FindUser(ctx, id)
	if err != nil {
		t.Fatalf("got error while finding user %v, expected nil", err)
	}

	if user.Name != someName {
		t.Fatalf("got user name %v, expected %v", user.Name, someName)
	}

	if user.ImageURL != someImageURL {
		t.Fatalf("got image url %v, expected %v", user.ImageURL, someImageURL)
	}
}

func TestFindUser_NotFound(t *testing.T) {
	db := newDB(t)
	defer db.close()

	ctx := t.Context()

	gotUser, err := storage.New(db, pkgslog.DiscardLogger).FindUser(ctx, 12345)
	if !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("got error while finding user %v, expected %v", err, user.ErrNotFound)
	}

	if gotUser.Name != "" {
		t.Fatalf("got user name %v, expected empty", gotUser.Name)
	}

	if gotUser.ImageURL != "" {
		t.Fatalf("got image url %v, expected empty", gotUser.ImageURL)
	}
}
