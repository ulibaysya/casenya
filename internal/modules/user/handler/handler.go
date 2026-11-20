package handler

import (
	"context"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/user/v1"
	"google.golang.org/grpc"
)

type Service interface {
	CreateUser(ctx context.Context, name string) (int64, error)
	GenerateUserOpaqueToken(ctx context.Context, id int64) (string, error)
	ExchangeUserOpaqueToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	srvc Service

	v1.UnimplementedUserServiceServer
	v1.UnimplementedBalanceServiceServer
	v1.UnimplementedItemServiceServer
	v1.UnimplementedOpaqueTokenServiceServer
}

func Register(s grpc.ServiceRegistrar, srvc Service) {
	h := &Handler{
		srvc: srvc,
	}
	_ = h
	v1.RegisterUserServiceServer(s, h)
	v1.RegisterBalanceServiceServer(s, h)
	v1.RegisterItemServiceServer(s, h)
	v1.RegisterOpaqueTokenServiceServer(s, h)
}
