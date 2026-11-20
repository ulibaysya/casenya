// TODO write FNX error test
package valkey

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/ulibaysya/casenya/internal/modules/user/userstorage"
)

func TestSaveUserSession(t *testing.T) {
	client := newValkey(t)
	resetValkey(t, client)
	defer client.Close()

	timeParam := time.Date(1961, time.April, 12, 6, 7, 0, 0, time.UTC)
	userStorage := userstorage.New(client, func() time.Time {
		return timeParam
	})

	ctx := t.Context()
	const userIDParam = 12345
	const sessionParam = "some-session"
	if err := userStorage.SaveUserSession(ctx, userIDParam, sessionParam); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	scan, err := client.Do(ctx, client.B().Scan().Cursor(0).Build()).AsScanEntry()
	if err != nil {
		t.Fatalf("error while performing SCAN: %v", err)
	}
	if scan.Cursor != 0 {
		t.Fatalf("expected scan cursor 0, got %v", scan.Cursor)
	}
	if elemsLen := len(scan.Elements); elemsLen != 1 {
		t.Fatalf("expected 1 element in scan, got %v", elemsLen)
	}
	const expectedKey = "session:" + sessionParam
	if scan.Elements[0] != expectedKey {
		t.Fatalf("expected key %v, got %v", expectedKey, scan.Elements[0])
	}

	sessionHash, err := client.Do(ctx, client.B().Hgetall().Key(expectedKey).Build()).AsStrMap()
	if err != nil {
		t.Fatalf("error while performing HGETALL: %v", err)
	}
	if sessionLen := len(sessionHash); sessionLen != 2 {
		t.Fatalf("expected 2 elements in session hash, got %v", sessionLen)
	}

	userIDValue, found := sessionHash["user_id"]
	if !found {
		t.Fatalf("not found user id in session hash")
	}
	expectedUserID := strconv.FormatInt(userIDParam, 10)
	if userIDValue != expectedUserID {
		t.Fatalf("expected user id %v, got %v", expectedUserID, userIDValue)
	}

	createdAtValue, found := sessionHash["created_at"]
	if !found {
		t.Fatalf("not found created at in session hash")
	}
	expectedCreatedAt := strconv.FormatInt(timeParam.Unix(), 10)
	if createdAtValue != expectedCreatedAt {
		t.Fatalf("expected created at %v, got %v", expectedCreatedAt, createdAtValue)
	}

	ttls, err := client.Do(ctx, client.B().Httl().Key(expectedKey).Fields().Numfields(2).Field("user_id").Field("created_at").Build()).AsIntSlice()
	if err != nil {
		t.Fatalf("error while performing HTTL: %v", err)
	}
	if ttlsLen := len(ttls); ttlsLen != 2 {
		t.Fatalf("expected 2 elements in ttls, got %v", ttlsLen)
	}
	if ttls[0] != ttls[1] {
		t.Fatalf("expected equal ttls, got non-equal: %v and %v", ttls[0], ttls[1])
	}
	const expectedTTL = 86400
	// using some leeway just in case
	if ttls[0] < expectedTTL-5 || ttls[0] > expectedTTL {
		t.Fatalf("expected ttl %v, got %v", expectedTTL, ttls[0])
	}
}

func TestFindUserSession(t *testing.T) {
	client := newValkey(t)
	resetValkey(t, client)
	defer client.Close()

	ctx := t.Context()

	const sessionParam = "some-session"
	const userIDValue = "123456"
	ret, err := client.Do(ctx, client.B().Hset().Key("session:"+sessionParam).FieldValue().FieldValue("user_id", userIDValue).Build()).AsInt64()
	if err != nil {
		t.Fatalf("error while inserting session: %v", err)
	}
	if ret != 1 {
		t.Fatalf("valkey returned not 1 when inserting session: %v", ret)
	}

	userStorage := userstorage.New(client, time.Now)

	userID, err := userStorage.FindUserSession(ctx, sessionParam)

	if err != nil {
		t.Fatalf("got %v error when finding user session, nil", err)
	}

	expectedUserID := mustParseInt(userIDValue, 10, 64)
	if userID != expectedUserID {
		t.Fatalf("got %v user id, expected %v", expectedUserID, userID)
	}
}

func TestFindUserSession_NotFound(t *testing.T) {
	client := newValkey(t)
	resetValkey(t, client)
	defer client.Close()

	userStorage := userstorage.New(client, time.Now)

	ctx := t.Context()
	const sessionParam = "some-session"
	userID, err := userStorage.FindUserSession(ctx, sessionParam)
	if !errors.Is(err, user.ErrSessionNotFound) {
		t.Fatalf("got %v error when finding user session, expected %v", err, user.ErrSessionNotFound)
	}
	if userID != 0 {
		t.Fatalf("got %v user id, expected 0", userID)
	}
}

func mustParseInt(s string, base int, bitSize int) int64 {
	res, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		panic(err)
	}
	return res
}
