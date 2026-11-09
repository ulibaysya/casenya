package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UploadItemImage(ctx context.Context, req *v1.UploadItemImageRequest) (*v1.UploadItemImageResponse, error) {
	err := validateUploadItemImageReq(req)
	if err != nil {
		return nil, err
	}

	uploadURL, err := h.srvc.UploadItemImage(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, item.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.UploadItemImageResponse{
		UploadUrl: uploadURL,
	}, nil
}

func validateUploadItemImageReq(req *v1.UploadItemImageRequest) error {
	const maxID = 1000000000
	if id := req.GetId(); id <= 0 || id > maxID {
		return status.Error(codes.InvalidArgument, "invalid id")
	}

	return nil
}
