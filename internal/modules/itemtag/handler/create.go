package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/itemtag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateTag(ctx context.Context, req *v1.CreateTagRequest) (*v1.CreateTagResponse, error) {
	err := validateCreateTagReq(req)
	if err != nil {
		return nil, err
	}

	if err = h.srvc.CreateTag(ctx, req.GetKey(), req.GetIsPublic(), req.GetValues()); err != nil {
		if errors.Is(err, itemtag.ErrKeyAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateTagResponse{}, nil
}

func validateCreateTagReq(req *v1.CreateTagRequest) error {
	key := req.GetKey()
	if keyLen := len(key); keyLen == 0 || keyLen > 100 {
		return status.Error(codes.InvalidArgument, "invalid key")
	}

	values := req.GetValues()
	if vLen := len(values); vLen == 0 || len(values) > 100 {
		return status.Error(codes.InvalidArgument, "invalid values")
	}
	valuesSet := map[string]struct{}{}
	for i, value := range values {
		if valueLen := len(value); valueLen == 0 || valueLen > 100 {
			return status.Errorf(codes.InvalidArgument, "invalid value %d", i)
		}
		if _, inSet := valuesSet[value]; inSet {
			return status.Errorf(codes.InvalidArgument, "value %d is specified more than once", i)
		}
		valuesSet[value] = struct{}{}
	}
	clear(valuesSet)

	return nil
}
