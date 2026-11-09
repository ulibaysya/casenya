package handler

import (
	"context"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"google.golang.org/grpc"
)

type Service interface {
	CreateUser(ctx context.Context, name string) (int64, error)
}

type Handler struct {
	srvc Service

	v1.UnimplementedUserServiceServer
}

func Register(s grpc.ServiceRegistrar, srvc Service) {
	v1.RegisterUserServiceServer(s, &Handler{
		srvc: srvc,
	})
}
