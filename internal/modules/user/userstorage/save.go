package userstorage

import (
	"context"
	"fmt"
	"strconv"
)

func (s *UserStorage) SaveUserSession(ctx context.Context, userID int64, session string) error {
	key := "session:" + session
	const ttl = 86400

	const userIDKey = "user_id"
	userIDValue := strconv.FormatInt(userID, 10)

	const createdAtKey = "created_at"
	createdAtValue := strconv.FormatInt(s.timeFunc().Unix(), 10)

	cmd := s.client.B().
		Hsetex().Key(key).Fnx().Ex(ttl).Fields().Numfields(2).FieldValue().
		FieldValue(userIDKey, userIDValue).
		FieldValue(createdAtKey, createdAtValue).
		Build()
	res, err := s.client.Do(ctx, cmd).AsInt64()
	if err != nil {
		return fmt.Errorf("saving user session with HSETEX: %w", err)
	}
	const hsetexSuccess = 1
	if res != hsetexSuccess {
		return fmt.Errorf("saving user session with HSETEX: unexpected result %d, expected %d (probably same session key with user_id and created_at specified, FNX trigger)", res, hsetexSuccess)
	}

	return nil
}
