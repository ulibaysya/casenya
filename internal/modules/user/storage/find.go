package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ulibaysya/casenya/internal/modules/user"
)

func (s *Storage) FindUser(ctx context.Context, id int64) (user.PublicUser, error) {
	u := user.PublicUser{}
	const userSelectSQL = `
		SELECT name, image_url
		FROM users
		WHERE id = $1;
	`
	imageURL := pgtype.Text{}
	err := s.db.QueryRow(ctx, userSelectSQL, id).Scan(&u.Name, &imageURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.PublicUser{}, user.ErrNotFound
		}
		return user.PublicUser{}, fmt.Errorf("selecting user: %w", err)
	}

	if imageURL.Valid {
		u.ImageURL = imageURL.String
	}

	return u, nil
}
