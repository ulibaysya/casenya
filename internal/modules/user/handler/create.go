package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/user/v1"
	"github.com/ulibaysya/casenya/internal/modules/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	err := validateCreateUserReq(req)
	if err != nil {
		return nil, err
	}

	id, err := h.srvc.CreateUser(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, user.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateUserResponse{
		Id: id,
	}, nil
}

func validateCreateUserReq(req *v1.CreateUserRequest) error {
	name := req.GetName()
	if nameLen := len(name); nameLen == 0 || nameLen > 100 {
		return status.Error(codes.InvalidArgument, "invalid name")
	}

	return nil
}
