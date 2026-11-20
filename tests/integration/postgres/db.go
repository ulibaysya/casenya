package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ulibaysya/casenya/internal/modules/user"
)

type testDB struct {
	*pgx.Conn
	t *testing.T
}

func newDB(t *testing.T) testDB {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	connstr, ok := os.LookupEnv("TEST_DB_CONNSTR")
	if !ok {
		t.Fatalf("TEST_DB_CONNSTR is not set for integration test")
	}

	db, err := pgx.Connect(ctx, connstr)
	if err != nil {
		t.Fatalf("error while initializing connection with db: %v", err)
	}

	if err = db.Ping(ctx); err != nil {
		t.Fatalf("error while pinging database: %v", err)
	}

	const resetSQL = `
		TRUNCATE
			items, 
			item_tag_keys, item_tag_values, item_tag_relations,
			items_prices,
			users, users_items,
			users_opaque_tokens
		RESTART IDENTITY CASCADE;
	`
	if _, err = db.Exec(ctx, resetSQL); err != nil {
		t.Fatalf("resetting db: %v", err)
	}

	return testDB{
		Conn: db,
		t:    t,
	}
}

func (d testDB) close() {
	d.t.Helper()

	ctx, cancel := context.WithTimeout(d.t.Context(), time.Second*5)
	defer cancel()

	if err := d.Close(ctx); err != nil {
		d.t.Logf("error while closing db connection: %v", err)
	}
}

func (d testDB) addUser(name string, imageURL string, balance int64) int64 {
	d.t.Helper()

	ctx, cancel := context.WithTimeout(d.t.Context(), time.Second*5)
	defer cancel()

	var id int64
	const insertSQL = `
		INSERT INTO users(name, image_url, balance)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err := d.QueryRow(ctx, insertSQL, name, imageURL, balance).Scan(&id)
	if err != nil {
		d.t.Fatalf("error while inserting user: %v", err)
	}

	return id
}

func (d testDB) addItem(name string) int64 {
	d.t.Helper()

	ctx, cancel := context.WithTimeout(d.t.Context(), time.Second*5)
	defer cancel()

	var id int64
	const insertSQL = `
		INSERT INTO items(name)
		VALUES ($1)
		RETURNING id
	`
	err := d.QueryRow(ctx, insertSQL, name).Scan(&id)
	if err != nil {
		d.t.Fatalf("error while inserting item: %v", err)
	}

	return id
}

func (d testDB) addUserItem(userID int64, itemID int64, attachedPrice int64, status user.ItemStatus) int64 {
	d.t.Helper()

	ctx, cancel := context.WithTimeout(d.t.Context(), time.Second*5)
	defer cancel()

	var id int64
	const insertSQL = `
		INSERT INTO users_items(
			user_id, 
			item_id, 
			attached_price, 
			status_id
		) 
		SELECT $1, $2, $3, s.id 
		FROM users_items_statuses AS s 
		WHERE s.code = $4 
		RETURNING id
	`
	err := d.QueryRow(ctx, insertSQL, userID, itemID, attachedPrice, status).Scan(&id)
	if err != nil {
		d.t.Fatalf("error while inserting item: %v", err)
	}

	return id
}
