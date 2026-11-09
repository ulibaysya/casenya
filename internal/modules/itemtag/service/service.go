package service

import (
	"context"
)

type Storage interface {
	CreateTag(ctx context.Context, name string, isPublic bool, values []string) error
}

type Service struct {
	strg Storage
}

func New(strg Storage) *Service {
	return &Service{
		strg: strg,
	}
}
