package storage

import (
	"context"
	"fmt"

	"github.com/ulibaysya/casenya/internal/modules/user"
)

func (s *Storage) SetUserTradeURL(ctx context.Context, id int64, tradeURL string) error {
	const tradeURLUpdateSQL = `
		UPDATE users
		SET trade_url = $1
		WHERE id = $2;
	`
	res, err := s.db.Exec(ctx, tradeURLUpdateSQL, tradeURL, id)
	if err != nil {
		return fmt.Errorf("updating user trade url: %w", err)
	}
	switch aff := res.RowsAffected(); aff {
	case 0:
		return user.ErrNotFound
	case 1:
	default:
		return fmt.Errorf("unexpected number of affected rows when updated user trade url: %v", aff)
	}

	return nil
}
