package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) FindPublicItem(ctx context.Context, req *v1.FindPublicItemRequest) (*v1.FindPublicItemResponse, error) {
	err := validateFindPublicItemReq(req)
	if err != nil {
		return nil, err
	}

	foundItem, err := h.srvc.FindPublicItem(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, item.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.FindPublicItemResponse{
		Name:     foundItem.Name,
		ImageUrl: foundItem.ImageURL,
		Tags:     makePublicReqTags(foundItem.Tags),
		Price:    foundItem.Price,
	}, nil
}

func validateFindPublicItemReq(req *v1.FindPublicItemRequest) error {
	const maxID = 1000000000
	if id := req.GetId(); id <= 0 || id > maxID {
		return status.Error(codes.InvalidArgument, "invalid id")
	}

	return nil
}

func makePublicReqTags(tags []item.Tag) []*v1.Tag {
	reqTags := make([]*v1.Tag, 0, len(tags))
	for _, tag := range tags {
		reqTags = append(reqTags, &v1.Tag{
			Key:    tag.Key,
			Values: tag.Values,
		})
	}

	return reqTags
}
