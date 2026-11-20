package userstorage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ulibaysya/casenya/internal/modules/user"
	"github.com/valkey-io/valkey-go"
)

func (s *UserStorage) FindUserSession(ctx context.Context, session string) (int64, error) {
	key := "session:" + session
	const userIDKey = "user_id"

	cmd := s.client.B().
		Hget().Key(key).Field(userIDKey).
		Build()
	userIDValue, err := s.client.Do(ctx, cmd).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return 0, user.ErrSessionNotFound
		}
		return 0, fmt.Errorf("finding user session with HGET: %w", err)
	}
	userID, err := strconv.ParseInt(userIDValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing user session value: %w", err)
	}

	return userID, nil
}
