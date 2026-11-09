package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"github.com/ulibaysya/casenya/internal/modules/itemtag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) FindItem(ctx context.Context, req *v1.FindItemRequest) (*v1.FindItemResponse, error) {
	err := validateFindItemReq(req)
	if err != nil {
		return nil, err
	}

	foundItem, err := h.srvc.FindItem(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, item.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.FindItemResponse{
		Name:     foundItem.Name,
		ImageUrl: foundItem.ImageURL,
		Tags:     makeReqTags(foundItem.Tags),
		Prices:   makeReqPrices(foundItem.Prices),
	}, nil
}

func validateFindItemReq(req *v1.FindItemRequest) error {
	const maxID = 1000000000
	if id := req.GetId(); id <= 0 || id > maxID {
		return status.Error(codes.InvalidArgument, "invalid id")
	}

	return nil
}

func makeReqPrices(prices []item.Price) []*v1.Price {
	reqPrices := make([]*v1.Price, 0, len(prices))
	for _, price := range prices {
		reqPrices = append(reqPrices, &v1.Price{
			Provider: string(price.Provider),
			Value:    price.Value,
		})
	}

	return reqPrices
}

func makeReqTags(tags []itemtag.Tag) []*v1.FullTag {
	reqTags := make([]*v1.FullTag, 0, len(tags))
	for _, tag := range tags {
		reqTags = append(reqTags, &v1.FullTag{
			Key:      tag.Key,
			IsPublic: tag.IsPublic,
			Values:   tag.Values,
		})
	}

	return reqTags
}
