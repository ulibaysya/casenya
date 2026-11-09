package service

import (
	"context"
	"log/slog"

	"github.com/ulibaysya/casenya/internal/modules/item"
)

type RNG interface {
	Int64() int64
}

type Storage interface {
	CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error)
	FindItem(ctx context.Context, id int64) (item.Item, error)
	SaveItemImageURL(ctx context.Context, id int64, url string) error
}

type BlobStorage interface {
	// The first string is a full URL where image should be uploaded.
	// The second string is a full URL where uploaded image will live.
	CreateItemImageUploadURL(ctx context.Context) (string, string, error)
	// // The second string is an object key where uploaded image will live.
	// GetUploadedItemImage(ctx context.Context, key string) (io.Reader, error)
}

type Service struct {
	logger   *slog.Logger
	strg     Storage
	blobStrg BlobStorage
	rng      RNG
}

func New(logger *slog.Logger, strg Storage, blobStrg BlobStorage, rng RNG) *Service {
	return &Service{
		logger:   logger,
		strg:     strg,
		blobStrg: blobStrg,
		rng:      rng,
	}
}
