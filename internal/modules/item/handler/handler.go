package handler

import (
	"context"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/v1"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"google.golang.org/grpc"
)

type Service interface {
	CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error)
	FindItem(ctx context.Context, id int64) (item.Item, error)
	FindPublicItem(ctx context.Context, id int64) (item.PublicItem, error)
	UploadItemImage(ctx context.Context, id int64) (string, error)
}

type Handler struct {
	srvc Service

	v1.UnimplementedItemServiceServer
}

func Register(s grpc.ServiceRegistrar, srvc Service) {
	v1.RegisterItemServiceServer(s, &Handler{
		srvc: srvc,
	})
}
