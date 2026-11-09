package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func newDB(t *testing.T) *pgx.Conn {
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

	return db
}

func resetDB(t *testing.T, db *pgx.Conn) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	const resetSQL = `
		TRUNCATE 
			items, 
			item_tag_keys,
			item_tag_values,
			item_tag_relations 
		RESTART IDENTITY CASCADE;
	`
	_, err := db.Exec(ctx, resetSQL)
	if err != nil {
		t.Fatalf("resetting db: %v", err)
	}
}

func closeDB(t *testing.T, db *pgx.Conn) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	if err := db.Close(ctx); err != nil {
		t.Logf("error while closing db connection: %v", err)
	}
}
