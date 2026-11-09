package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateItem(ctx context.Context, req *v1.CreateItemRequest) (*v1.CreateItemResponse, error) {
	err := validateCreateItemReq(req)
	if err != nil {
		return nil, err
	}

	id, err := h.srvc.CreateItem(ctx, req.GetName(), req.GetPrice(), makeTags(req.GetTags()))
	if err != nil {
		if errors.Is(err, item.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateItemResponse{
		Id: id,
	}, nil
}

func validateCreateItemReq(req *v1.CreateItemRequest) error {
	name := req.GetName()
	if nameLen := len(name); nameLen == 0 || nameLen > 100 {
		return status.Error(codes.InvalidArgument, "invalid name")
	}

	price := req.GetPrice()
	const maxPrice = 1000000000
	if price <= 0 || price > maxPrice {
		return status.Error(codes.InvalidArgument, "invalid price")
	}

	tags := req.GetTags()
	keysSet := make(map[string]struct{}, len(tags))
	for i, tag := range tags {
		tagKey := tag.GetKey()
		if keyLen := len(tagKey); keyLen == 0 {
			return status.Errorf(codes.InvalidArgument, "key %d has invalid name", i)
		}
		if _, inSet := keysSet[tag.Key]; inSet {
			return status.Errorf(codes.InvalidArgument, "key %d is specified more than once", i)
		}
		keysSet[tagKey] = struct{}{}

		tagValues := tag.GetValues()
		if valuesLen := len(tagValues); valuesLen == 0 || valuesLen > 100 {
			return status.Errorf(codes.InvalidArgument, "key %d has invalid number of values", i)
		}
		valuesSet := make(map[string]struct{}, len(tag.Values))
		for j, value := range tagValues {
			if valueLen := len(value); valueLen == 0 || valueLen > 100 {
				return status.Errorf(codes.InvalidArgument, "value %d of key %d has invalid name", j, i)
			}
			if _, inSet := valuesSet[value]; inSet {
				return status.Errorf(codes.InvalidArgument, "value %d of key %d is specified more than once", j, i)
			}
			valuesSet[value] = struct{}{}
		}
		clear(valuesSet)
	}
	clear(keysSet)

	return nil
}

func makeTags(reqTags []*v1.Tag) []item.Tag {
	tags := make([]item.Tag, 0, len(reqTags))
	for _, reqTag := range reqTags {
		tags = append(tags, item.Tag{
			Key:    reqTag.Key,
			Values: reqTag.GetValues(),
		})
	}
	return tags
}
