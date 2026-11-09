package integration

import (
	"testing"

	"github.com/ulibaysya/casenya/internal/modules/itemtag/storage"
	pkgslog "github.com/ulibaysya/casenya/pkg/slog"
)

func TestItemtagCreate(t *testing.T) {
	db := newDB(t)
	resetDB(t, db)
	defer closeDB(t, db)

	err := storage.New(pkgslog.DiscardLogger, db).CreateKey(t.Context(), "hi", true, []string{"val1", "val2"})
	if err != nil {
		t.Fatal(err)
	}
}
