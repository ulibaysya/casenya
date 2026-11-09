package service

import (
	"context"

	"github.com/ulibaysya/casenya/internal/modules/item"
)


func (s *Service) FindItem(ctx context.Context, id int64) (item.Item, error) {
	return s.strg.FindItem(ctx, id)
}
