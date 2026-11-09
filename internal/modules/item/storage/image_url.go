package storage

import (
	"context"
	"fmt"
)

func (s *Storage) SaveItemImageURL(ctx context.Context, id int64, url string) error {
	const updateSQL = `
		UPDATE items
		SET image_url = $1
		WHERE id = $2
	`
	if _, err := s.db.Exec(ctx, updateSQL, url, id); err != nil {
		return fmt.Errorf("updating image url: %w", err)
	}

	return nil
}
