package service

import (
	"context"
	"slices"

	"github.com/ulibaysya/casenya/internal/modules/item"
)

func (s *Service) FindPublicItem(ctx context.Context, id int64) (item.PublicItem, error) {
	foundItem, err := s.strg.FindItem(ctx, id)
	if err != nil {
		return item.PublicItem{}, err
	}

	tags := []item.Tag{}
	for _, tag := range foundItem.Tags {
		if tag.IsPublic {
			tags = append(tags, item.Tag{
				Key:    tag.Key,
				Values: tag.Values,
			})
		}
	}

	var price int64
	if len(foundItem.Prices) != 0 {
		prices := make([]int64, 0, len(foundItem.Prices))
		for _, price := range foundItem.Prices {
			prices = append(prices, price.Value)
		}
		price = slices.Min(prices)
	}

	return item.PublicItem{
		Name:     foundItem.Name,
		ImageURL: foundItem.ImageURL,
		Price:    price,
		Tags:     tags,
	}, nil
}
