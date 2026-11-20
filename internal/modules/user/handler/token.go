package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/user/v1"
	"github.com/ulibaysya/casenya/internal/modules/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GenerateUserOpaqueToken(ctx context.Context, req *v1.GenerateUserOpaqueTokenRequest) (*v1.GenerateUserOpaqueTokenResponse, error) {
	err := validateGenerateUserOpaqueTokenRequest(req)
	if err != nil {
		return nil, err
	}

	token, err := h.srvc.GenerateUserOpaqueToken(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.GenerateUserOpaqueTokenResponse{
		Token: token,
	}, nil
}

func validateGenerateUserOpaqueTokenRequest(req *v1.GenerateUserOpaqueTokenRequest) error {
	id := req.GetId()
	if id <= 0 {
		return status.Error(codes.InvalidArgument, "invalid id")
	}

	return nil
}

func (h *Handler) ExchangeUserOpaqueToken(ctx context.Context, req *v1.ExchangeUserOpaqueTokenRequest) (*v1.ExchangeUserOpaqueTokenResponse, error) {
	err := validateExchangeUserOpaqueTokenRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := h.srvc.ExchangeUserOpaqueToken(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, user.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.ExchangeUserOpaqueTokenResponse{
		Id: id,
	}, nil
}

func validateExchangeUserOpaqueTokenRequest(req *v1.ExchangeUserOpaqueTokenRequest) error {
	token := req.GetToken()
	if len(token) == 0 {
		return status.Error(codes.InvalidArgument, "invalid token")
	}

	return nil
}
